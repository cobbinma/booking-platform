use crate::models;
use async_trait::async_trait;
use chrono::{DateTime, Datelike, Duration, NaiveDate, NaiveTime, TimeZone, Timelike, Utc};
use mockall::*;
use protobuf::booking::api::booking_api_server::BookingApi;
use protobuf::booking::api::GetSlotResponse;
use protobuf::booking::models::{Booking, Slot, SlotInput};
use protobuf::venue::models::Venue;
use std::collections::HashSet;
use std::ops::Add;
use tonic::{Request, Response, Status};
use uuid::Uuid;

#[automock]
#[async_trait]
pub trait VenueClient {
    async fn get_venue(&self, venue_id: String) -> Result<Venue, Status>;
}

#[automock]
#[async_trait]
pub trait TableClient {
    async fn get_tables_with_capacity(
        &self,
        venue_id: String,
        capacity: u32,
    ) -> Result<Vec<String>, Status>;
}

#[automock]
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
        tracing::info!("get slot call for venue '{}', starting at '{}', duration '{}' minutes, for '{} people'", &slot.venue_id, &slot.starts_at, slot.duration, slot.people);

        let slot_starts_at = DateTime::parse_from_rfc3339(&slot.starts_at).map_err(|e| {
            log::error!("could not parse starting date time : {}", e);
            Status::invalid_argument("could not parse starting date time")
        })?;

        let slot_date = NaiveDate::from_ymd(
            slot_starts_at.year(),
            slot_starts_at.month(),
            slot_starts_at.day(),
        );

        let (opens, closes) = &self
            .get_opening_times(slot.venue_id.clone(), slot_date)
            .await?;

        if slot_starts_at < *opens
            || slot_starts_at + Duration::minutes(slot.duration as i64) > *closes
        {
            return Err(Status::invalid_argument("venue is closed at that time"));
        }

        let tables_with_capacity = &self
            .table_client
            .get_tables_with_capacity(slot.venue_id.clone(), slot.people)
            .await?;

        if tables_with_capacity.is_empty() {
            return Err(Status::invalid_argument(
                "venue does not have tables that large",
            ));
        }

        let bookings = self.get_bookings_by_date(&slot.venue_id.clone(), &slot_date)?;

        let mut free_time_slots = HashSet::new();
        // Loop through all possible time slots for the desired date to see which slots have a table free
        let mut slot_time = *opens;
        while slot_time <= *closes - Duration::minutes(slot.duration as i64) {
            let free_table_id = get_free_table(
                slot.duration as i64,
                tables_with_capacity,
                &bookings,
                &slot_time,
            );

            if free_table_id.is_some() {
                free_time_slots.insert(slot_time);
            }

            slot_time = slot_time + Duration::minutes(30);
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
        tracing::info!("create booking call for venue '{}', starting at '{}', duration '{}' minutes, for '{} people'", &slot.venue_id, &slot.starts_at, slot.duration, slot.people);

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
            .get_opening_times(slot.venue_id.clone(), slot_date)
            .await?;

        if slot_starts_at < *opens
            || slot_starts_at + Duration::minutes(slot.duration as i64) > *closes
        {
            return Err(Status::invalid_argument("venue is closed at that time"));
        }

        let tables_with_capacity = &self
            .table_client
            .get_tables_with_capacity(slot.venue_id.clone(), slot.people)
            .await?;

        if tables_with_capacity.is_empty() {
            return Err(Status::invalid_argument(
                "venue does not have tables that large",
            ));
        }

        let bookings = self.get_bookings_by_date(&slot.venue_id.clone(), &slot_date)?;

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
                .filter(|&booking| booking.table_id.to_string() == **table_id)
                .all(|b| {
                    !(*starts_at < b.ends_at
                        && b.starts_at < *starts_at + Duration::minutes(duration))
                })
        })
        .map(|id| id.clone())
        .next()
}

impl BookingService {
    fn get_bookings_by_date(
        &self,
        venue_id: &str,
        day: &NaiveDate,
    ) -> Result<Vec<models::Booking>, Status> {
        Ok(self
            .repository
            .get_bookings_by_date(
                &Uuid::parse_str(venue_id).map_err(|e| {
                    log::error!("could not parse uuid : {}", e);
                    Status::internal("could not parse uuid")
                })?,
                &day,
            )?
            .iter()
            .map(|booking| booking.clone())
            .collect::<Vec<models::Booking>>())
    }

    async fn get_opening_times(
        &self,
        venue_id: String,
        date: NaiveDate,
    ) -> Result<(DateTime<Utc>, DateTime<Utc>), Status> {
        tracing::debug!(
            "getting opening times for venue {} on date {}",
            &venue_id,
            &date.to_string()
        );

        let venue = &self.venue_client.get_venue(venue_id).await?;

        let opening_hours_specification = venue
            .opening_hours
            .iter()
            .filter(|&hours| hours.day_of_week == date.weekday().number_from_monday())
            .next()
            .ok_or_else(|| Status::invalid_argument("venue not open on given date"))?;

        fn combine_date_and_time(date: NaiveDate, c: NaiveTime) -> DateTime<Utc> {
            Utc.ymd(date.year(), date.month(), date.day())
                .and_hms(c.hour(), c.minute(), c.second())
        }

        let opens = NaiveTime::parse_from_str(&opening_hours_specification.opens, "%H:%M")
            .map_err(|e| {
                log::error!("could not parse opens time : {}", e);
                Status::internal("could not parse opens time")
            })
            .map(|o| combine_date_and_time(date, o))?;

        let closes = NaiveTime::parse_from_str(&opening_hours_specification.closes, "%H:%M")
            .map_err(|e| {
                log::error!("could not parse closes time : {}", e);
                Status::internal("could not parse closes time")
            })
            .map(|c| combine_date_and_time(date, c))?;

        Ok((opens, closes))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::models::Booking;
    use chrono::{DateTime, Duration, NaiveDateTime, Utc};
    use protobuf::venue::models::OpeningHoursSpecification;
    use uuid::Uuid;

    #[test]
    fn test_get_free_table_no_bookings() {
        let free_table = super::get_free_table(
            60,
            &vec!["3a3789ca-7174-4127-ae50-a644d69f1d27".to_string()],
            &vec![],
            &DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(61, 0), Utc),
        );

        assert_eq!(
            free_table,
            Some("3a3789ca-7174-4127-ae50-a644d69f1d27".to_string())
        )
    }

    #[test]
    fn test_get_free_table_with_booking_conflict() {
        let starts_at = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(61, 0), Utc);
        let free_table = super::get_free_table(
            60,
            &vec!["3a3789ca-7174-4127-ae50-a644d69f1d27".to_string()],
            &vec![Booking {
                id: Default::default(),
                customer_email: "".to_string(),
                venue_id: Default::default(),
                table_id: Uuid::parse_str("3a3789ca-7174-4127-ae50-a644d69f1d27")
                    .expect("could not parse uuid"),
                people: 4,
                date: starts_at.date().naive_utc(),
                starts_at,
                ends_at: starts_at + Duration::minutes(60),
                duration: 60,
            }],
            &starts_at,
        );

        assert_eq!(free_table, None)
    }

    #[test]
    fn test_get_free_table_with_partial_booking_conflict() {
        let starts_at = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(61, 0), Utc);
        let free_table = super::get_free_table(
            60,
            &vec!["3a3789ca-7174-4127-ae50-a644d69f1d27".to_string()],
            &vec![Booking {
                id: Default::default(),
                customer_email: "".to_string(),
                venue_id: Default::default(),
                table_id: Uuid::parse_str("3a3789ca-7174-4127-ae50-a644d69f1d27")
                    .expect("could not parse uuid"),
                people: 4,
                date: starts_at.date().naive_utc(),
                starts_at: starts_at + Duration::minutes(30),
                ends_at: starts_at + Duration::minutes(90),
                duration: 60,
            }],
            &starts_at,
        );

        assert_eq!(free_table, None)
    }

    #[test]
    fn test_get_free_table_with_no_tables() {
        let starts_at = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(61, 0), Utc);
        let free_table = super::get_free_table(60, &vec![], &vec![], &starts_at);

        assert_eq!(free_table, None)
    }

    #[tokio::test]
    async fn test_get_opening_times() {
        let date = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704678400, 0), Utc);
        let mut mock = MockVenueClient::new();
        mock.expect_get_venue()
            .with(predicate::eq(
                "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
            ))
            .times(1)
            .returning(|_| {
                Ok(Venue {
                    id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                    name: "".to_string(),
                    opening_hours: vec![OpeningHoursSpecification {
                        day_of_week: 5,
                        opens: "10:00".to_string(),
                        closes: "22:00".to_string(),
                        valid_from: "".to_string(),
                        valid_through: "".to_string(),
                    }],
                    special_opening_hours: vec![],
                })
            });
        let service = BookingService::new(
            Box::new(MockRepository::new()),
            Box::new(mock),
            Box::new(MockTableClient::new()),
        )
        .expect("could not construct booking service");

        let result = service
            .get_opening_times(
                "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                date.naive_utc().date(),
            )
            .await
            .expect("did not expect error");

        assert_eq!(
            result,
            (
                DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704714400, 0), Utc),
                DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704757600, 0), Utc)
            )
        )
    }
}
