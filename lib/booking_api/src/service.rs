use crate::models;
use crate::postgres::Postgres;
use crate::schema::bookings::columns::venue_id;
use async_trait::async_trait;
use chrono::format::Numeric::Timestamp;
use chrono::{DateTime, Datelike, Duration, NaiveDate, NaiveTime, TimeZone, Timelike, Utc};
use protobuf::booking::api::booking_api_server::BookingApi;
use protobuf::booking::api::GetSlotResponse;
use protobuf::booking::models::{Booking, Slot, SlotInput};
use protobuf::venue::api::table_api_client::TableApiClient;
use protobuf::venue::api::table_api_server::TableApi;
use protobuf::venue::api::venue_api_client::VenueApiClient;
use protobuf::venue::api::venue_api_server::VenueApi;
use protobuf::venue::api::{GetTablesRequest, GetVenueRequest};
use std::borrow::Borrow;
use std::collections::{HashMap, HashSet};
use std::ops::{Add, Deref};
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

#[async_trait]
pub trait VenueClient {
    async fn get_opening_times(
        &self,
        venue_id: String,
        date: NaiveDate,
    ) -> Result<(DateTime<Utc>, DateTime<Utc>), Status>;
}

#[async_trait]
pub trait TableClient {
    async fn get_tables_with_capacity(
        &self,
        venue_id: String,
        capacity: u32,
    ) -> Result<Vec<String>, Status>;
}

pub trait Repository {
    fn get_bookings_by_date(
        &self,
        venue: &Uuid,
        day: &NaiveDate,
    ) -> Result<Vec<models::Booking>, Status>;

    fn create_booking(&self, new_booking: &models::Booking) -> Result<(), Status>;
}

pub struct BookingService {
    repository: Box<dyn Repository + Send + Sync + 'static>,
    venue_client: Box<dyn VenueClient + Send + Sync + 'static>,
    table_client: Box<dyn TableClient + Send + Sync + 'static>,
}

impl BookingService {
    pub fn new(
        repository: Box<dyn Repository + Send + Sync + 'static>,
        venue_client: Box<dyn VenueClient + Send + Sync + 'static>,
        table_client: Box<dyn TableClient + Send + Sync + 'static>,
    ) -> Result<Self, Box<dyn std::error::Error>> {
        Ok(BookingService {
            repository,
            venue_client,
            table_client,
        })
    }
}

#[async_trait]
impl BookingApi for BookingService {
    async fn get_slot(&self, req: Request<SlotInput>) -> Result<Response<GetSlotResponse>, Status> {
        let slot = req.into_inner();

        let slot_starts_at = DateTime::parse_from_rfc3339(&slot.starts_at).map_err(|e| {
            log::error!("could not parse date : {}", e);
            Status::internal("could not parse date")
        })?;
        let slot_date = NaiveDate::from_ymd(
            slot_starts_at.year(),
            slot_starts_at.month(),
            slot_starts_at.day(),
        );

        let (opens, closes) = &self
            .venue_client
            .get_opening_times(slot.venue_id.clone(), slot_date)
            .await?;

        if slot_starts_at < *opens
            || slot_starts_at + Duration::minutes(slot.duration as i64) > *closes
        {
            return Err(Status::invalid_argument("venue is closed"));
        }

        let tables_with_capacity = &self
            .table_client
            .get_tables_with_capacity(slot.venue_id.clone(), slot.people)
            .await?;

        if tables_with_capacity.is_empty() {
            return Err(Status::invalid_argument(
                "restaurant does not have tables that large",
            ));
        }

        let bookings = self.get_bookings_by_date(&slot, &slot_date)?;

        let mut free_time_slots = HashSet::new();
        let mut t = *opens;
        while t <= *closes - Duration::minutes(slot.duration as i64) {
            let free_table_id =
                get_free_table(slot.duration as i64, tables_with_capacity, &bookings, &t);

            if free_table_id.is_some() {
                free_time_slots.insert(t);
            }

            t = t + Duration::minutes(30);
        }

        let other_available_slots: Vec<Slot> = free_time_slots
            .iter()
            .map(|(time)| Slot {
                venue_id: slot.venue_id.clone(),
                email: slot.email.clone(),
                people: slot.people,
                starts_at: time.to_rfc3339(),
                ends_at: (*time + Duration::minutes(slot.duration as i64)).to_rfc3339(),
                duration: slot.duration,
            })
            .collect();

        Ok(Response::new(GetSlotResponse {
            r#match: free_time_slots
                .get(&slot_starts_at.with_timezone(&Utc))
                .map(|_| Slot {
                    venue_id: slot.venue_id,
                    email: slot.email,
                    people: slot.people,
                    starts_at: slot.starts_at,
                    ends_at: (slot_starts_at + Duration::minutes(slot.duration as i64))
                        .to_rfc3339(),
                    duration: slot.duration,
                }),
            other_available_slots,
        }))
    }

    async fn create_booking(&self, req: Request<SlotInput>) -> Result<Response<Booking>, Status> {
        let slot = req.into_inner();

        let slot_starts_at = DateTime::parse_from_rfc3339(&slot.starts_at).map_err(|e| {
            log::error!("could not parse date : {}", e);
            Status::internal("could not parse date")
        })?;
        let slot_date = NaiveDate::from_ymd(
            slot_starts_at.year(),
            slot_starts_at.month(),
            slot_starts_at.day(),
        );

        let (opens, closes) = &self
            .venue_client
            .get_opening_times(slot.venue_id.clone(), slot_date)
            .await?;

        if slot_starts_at < *opens
            || slot_starts_at + Duration::minutes(slot.duration as i64) > *closes
        {
            return Err(Status::invalid_argument("venue is closed"));
        }

        let tables_with_capacity = &self
            .table_client
            .get_tables_with_capacity(slot.venue_id.clone(), slot.people)
            .await?;

        if tables_with_capacity.is_empty() {
            return Err(Status::invalid_argument(
                "restaurant does not have tables that large",
            ));
        }

        let bookings = self.get_bookings_by_date(&slot, &slot_date)?;

        let free_table_id = get_free_table(
            slot.duration as i64,
            tables_with_capacity,
            &bookings,
            &slot_starts_at.with_timezone(&Utc),
        );

        if let Some(table_id) = free_table_id {
            let id = uuid::Uuid::new_v4();
            let new_booking = models::Booking {
                id: id.clone(),
                customer_email: slot.email.clone(),
                venue_id: Uuid::parse_str(&slot.venue_id)
                    .map_err(|_| Status::invalid_argument("could not parse uuid"))?,
                table_id: Uuid::parse_str(&table_id)
                    .map_err(|_| Status::internal("could not parse table uuid"))?,
                people: slot.people as i32,
                date: slot_date,
                starts_at: slot_starts_at.with_timezone(&Utc),
                ends_at: slot_starts_at
                    .with_timezone(&Utc)
                    .add(Duration::minutes(slot.duration as i64)),
                duration: slot.duration as i32,
            };

            self.repository.create_booking(&new_booking)?;

            Ok(Response::new(Booking {
                id: id.to_string(),
                venue_id: slot.venue_id.clone(),
                email: slot.email.clone(),
                people: slot.people,
                starts_at: slot.starts_at,
                ends_at: (slot_starts_at + Duration::minutes(slot.duration as i64)).to_rfc3339(),
                duration: slot.duration,
                table_id: table_id.to_string(),
            }))
        } else {
            Err(Status::not_found("could not find a free slot"))
        }
    }
}

fn get_free_table(
    duration: i64,
    tables_with_capacity: &Vec<String>,
    bookings: &Vec<models::Booking>,
    starts_at: &DateTime<Utc>,
) -> Option<String> {
    tables_with_capacity
        .iter()
        .filter(|table_id| {
            bookings
                .iter()
                .filter(|booking| booking.table_id.to_string() == **table_id)
                .all(|b| {
                    !(*starts_at < b.starts_at
                        && b.starts_at < *starts_at + Duration::minutes(duration))
                })
        })
        .map(|id| id.clone())
        .next()
}

impl BookingService {
    fn get_bookings_by_date(
        &self,
        slot: &SlotInput,
        day: &NaiveDate,
    ) -> Result<Vec<models::Booking>, Status> {
        Ok(self
            .repository
            .get_bookings_by_date(
                &Uuid::parse_str(&slot.venue_id).map_err(|e| {
                    log::error!("could not parse uuid : {}", e);
                    Status::internal("could not parse uuid")
                })?,
                &day,
            )?
            .iter()
            .map(|booking| booking.clone())
            .collect::<Vec<models::Booking>>())
    }
}
