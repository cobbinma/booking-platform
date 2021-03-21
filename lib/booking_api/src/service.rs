use crate::models;
use async_trait::async_trait;
use chrono::{DateTime, Datelike, Duration, NaiveDate, NaiveTime, TimeZone, Timelike, Utc};
use num::integer::Integer;
use protobuf::booking::api::booking_api_server::BookingApi;
use protobuf::booking::api::{CancelBookingRequest, GetBookingsRequest, GetBookingsResponse, GetSlotResponse, SlotInput, BookingInput};
use protobuf::booking::models::{Booking, Slot};
use protobuf::venue::models::Venue;
use std::collections::HashSet;
use std::ops::Add;
use tonic::{Request, Response, Status};
use uuid::Uuid;

#[cfg(test)]
use mockall::automock;

#[cfg_attr(test, automock)]
#[async_trait]
pub trait VenueClient {
    async fn get_venue(&self, venue_id: String) -> Result<Venue, Status>;
    async fn get_tables_with_capacity(
        &self,
        venue_id: String,
        capacity: u32,
    ) -> Result<Vec<String>, Status>;
}

#[derive(PartialEq, Debug)]
pub struct BookingsFilter {
    pub venue: Option<Uuid>,
    pub day: Option<NaiveDate>,
}

#[cfg_attr(test, automock)]
pub trait Repository {
    fn get_bookings(
        &self,
        filter: BookingsFilter,
        page: Option<i32>,
        limit: Option<i32>,
    ) -> Result<Vec<models::Booking>, Status>;
    fn create_booking(&self, new_booking: &models::Booking) -> Result<(), Status>;
    fn cancel_booking(&self, id: &Uuid) -> Result<models::Booking, Status>;
    fn count_bookings(&self, filter: &BookingsFilter) -> Result<i64, Status>;
}

pub struct BookingService {
    repository: Box<dyn Repository + Send + Sync + 'static>,
    venue_client: Box<dyn VenueClient + Send + Sync + 'static>,
    uuid: Box<dyn UuidGetter + Send + Sync + 'static>,
}

impl BookingService {
    pub fn new(
        repository: Box<dyn Repository + Send + Sync + 'static>,
        venue_client: Box<dyn VenueClient + Send + Sync + 'static>,
        custom_uuid: Option<Box<dyn UuidGetter + Send + Sync + 'static>>,
    ) -> Result<Self, Box<dyn std::error::Error>> {
        let uuid = custom_uuid.unwrap_or_else(|| Box::new(GetUuid::default()));

        Ok(BookingService {
            repository,
            venue_client,
            uuid,
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

        let ((opens, closes), mut tables_with_capacity) = tokio::try_join!(
            self.get_opening_times(slot.venue_id.clone(), slot_date),
            self.venue_client
                .get_tables_with_capacity(slot.venue_id.clone(), slot.people)
        )?;

        if slot_starts_at < opens
            || slot_starts_at + Duration::minutes(slot.duration as i64) > closes
        {
            return Err(Status::invalid_argument("venue is closed at that time"));
        }

        if tables_with_capacity.is_empty() {
            return Err(Status::invalid_argument(
                "venue does not have tables that large",
            ));
        }

        let bookings = self.get_bookings_by_date(&slot.venue_id.clone(), &slot_date)?;

        let mut free_time_slots = HashSet::new();
        // Loop through all possible time slots for the desired date to see which slots have a table free
        let mut slot_time = opens;
        while slot_time <= closes - Duration::minutes(slot.duration as i64) {
            let free_table_id = get_free_table(
                slot.duration as i64,
                &mut tables_with_capacity,
                &bookings,
                &slot_time,
            );

            if free_table_id.is_some() {
                free_time_slots.insert(slot_time);
            }

            slot_time = slot_time + Duration::minutes(30);
        }

        let mut other_available_times: Vec<&DateTime<Utc>> =
            free_time_slots.iter().collect::<Vec<&DateTime<Utc>>>();
        other_available_times.sort();

        let other_available_slots: Vec<Slot> = other_available_times
            .iter()
            .map(|time| Slot {
                venue_id: slot.venue_id.clone(),
                email: slot.email.clone(),
                people: slot.people,
                starts_at: time.to_rfc3339(),
                ends_at: (**time + Duration::minutes(slot.duration as i64)).to_rfc3339(),
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

    async fn create_booking(&self, req: Request<BookingInput>) -> Result<Response<Booking>, Status> {
        let input = req.into_inner();
        tracing::info!("create booking call for venue '{}', starting at '{}', duration '{}' minutes, for '{} people'", &input.venue_id, &input.starts_at, input.duration, input.people);

        let slot_starts_at = DateTime::parse_from_rfc3339(&input.starts_at).map_err(|e| {
            log::error!("could not parse date : {}", e);
            Status::internal("could not parse date")
        })?;

        let slot_date = NaiveDate::from_ymd(
            slot_starts_at.year(),
            slot_starts_at.month(),
            slot_starts_at.day(),
        );

        let ((opens, closes), mut tables_with_capacity) = tokio::try_join!(
            self.get_opening_times(input.venue_id.clone(), slot_date),
            self.venue_client
                .get_tables_with_capacity(input.venue_id.clone(), input.people)
        )?;

        if slot_starts_at < opens
            || slot_starts_at + Duration::minutes(input.duration as i64) > closes
        {
            return Err(Status::invalid_argument("venue is closed at that time"));
        }

        if tables_with_capacity.is_empty() {
            return Err(Status::invalid_argument(
                "venue does not have tables that large",
            ));
        }

        let bookings = self.get_bookings_by_date(&input.venue_id.clone(), &slot_date)?;

        let free_table_id = get_free_table(
            input.duration as i64,
            &mut tables_with_capacity,
            &bookings,
            &slot_starts_at.with_timezone(&Utc),
        );

        if let Some(table_id) = free_table_id {
            let id = self.uuid.uuid();
            let new_booking = models::Booking {
                id,
                customer_email: input.email.clone(),
                venue_id: Uuid::parse_str(&input.venue_id)
                    .map_err(|_| Status::invalid_argument("could not parse uuid"))?,
                table_id: Uuid::parse_str(&table_id)
                    .map_err(|_| Status::internal("could not parse table uuid"))?,
                people: input.people as i32,
                date: slot_date,
                starts_at: slot_starts_at.with_timezone(&Utc),
                ends_at: slot_starts_at
                    .with_timezone(&Utc)
                    .add(Duration::minutes(input.duration as i64)),
                duration: input.duration as i32,
                name: input.name.clone(),
                given_name: input.given_name.clone()
            };

            self.repository.create_booking(&new_booking)?;

            Ok(Response::new(Booking {
                id: id.to_string(),
                venue_id: input.venue_id.clone(),
                email: input.email.clone(),
                people: input.people,
                starts_at: input.starts_at,
                ends_at: (slot_starts_at + Duration::minutes(input.duration as i64)).to_rfc3339(),
                duration: input.duration,
                table_id,
                name: input.name,
                given_name: input.given_name
            }))
        } else {
            Err(Status::not_found("could not find a free slot"))
        }
    }

    async fn get_bookings(
        &self,
        request: Request<GetBookingsRequest>,
    ) -> Result<Response<GetBookingsResponse>, Status> {
        let req = request.into_inner();
        tracing::info!(
            "get bookings call for venue '{}', date '{}'",
            &req.venue_id,
            &req.date
        );

        let venue = match req.venue_id.is_empty() {
            true => None,
            false => Some(Uuid::parse_str(&req.venue_id).map_err(|e| {
                log::error!("could not parse uuid : {}", e);
                Status::invalid_argument("could not parse uuid")
            })?),
        };

        let day = match req.date.is_empty() {
            true => None,
            false => Some(
                DateTime::parse_from_rfc3339(&req.date)
                    .map(|dt| NaiveDate::from_ymd(dt.year(), dt.month(), dt.day()))
                    .map_err(|e| {
                        log::error!("could not parse date : {}", e);
                        Status::invalid_argument("could not parse date")
                    })?,
            ),
        };

        let filter = BookingsFilter { venue, day };

        let count = self.repository.count_bookings(&filter)?;
        tracing::debug!("found {} bookings", count);

        let mut bookings = self
            .repository
            .get_bookings(filter, Some(req.page), Some(req.limit))?;
        tracing::debug!("got {} bookings", bookings.len());

        let has_next_page = bookings.len() >= req.limit as usize;

        if bookings.len() > req.limit as usize {
            tracing::debug!("removing extra booking");
            bookings.pop();
        }

        Ok(Response::new(GetBookingsResponse {
            bookings: bookings
                .iter()
                .map(|b| Booking {
                    id: b.id.to_string(),
                    venue_id: b.venue_id.to_string(),
                    email: b.customer_email.clone(),
                    people: b.people as u32,
                    starts_at: b.starts_at.to_rfc3339(),
                    ends_at: b.ends_at.to_rfc3339(),
                    duration: b.duration as u32,
                    table_id: b.table_id.to_string(),
                    name: b.name.clone(),
                    given_name: b.given_name.clone(),
                })
                .collect(),
            has_next_page,
            pages: (count as i32).div_ceil(&req.limit),
        }))
    }

    async fn cancel_booking(
        &self,
        request: Request<CancelBookingRequest>,
    ) -> Result<Response<Booking>, Status> {
        let req = request.into_inner();
        tracing::info!("cancel booking call for id '{}'", &req.id);

        let removed = self
            .repository
            .cancel_booking(&Uuid::parse_str(&req.id).map_err(|e| {
                log::error!("could not parse uuid : {}", e);
                Status::invalid_argument("could not parse uuid")
            })?)
            .map(|b| Booking {
                id: b.id.to_string(),
                venue_id: b.venue_id.to_string(),
                email: b.customer_email.clone(),
                people: b.people as u32,
                starts_at: b.starts_at.to_rfc3339(),
                ends_at: b.ends_at.to_rfc3339(),
                duration: b.duration as u32,
                table_id: b.table_id.to_string(),
                name: b.name.clone(),
                given_name: b.given_name.clone()
            })?;

        Ok(Response::new(removed))
    }
}

fn get_free_table(
    duration: i64,
    tables_with_capacity: &mut [String],
    bookings: &[models::Booking],
    starts_at: &DateTime<Utc>,
) -> Option<String> {
    tables_with_capacity.sort_by(|a, b| {
        bookings
            .iter()
            .filter(|&booking| a == &booking.table_id.to_string())
            .count()
            .cmp(
                &bookings
                    .iter()
                    .filter(|&booking| b == &booking.table_id.to_string())
                    .count(),
            )
    });
    tables_with_capacity
        .iter()
        .find(|table_id| {
            bookings
                .iter()
                .filter(|&booking| booking.table_id.to_string() == **table_id)
                .all(|b| {
                    !(*starts_at < b.ends_at
                        && b.starts_at < *starts_at + Duration::minutes(duration))
                })
        })
        .cloned()
}

impl BookingService {
    fn get_bookings_by_date(
        &self,
        venue_id: &str,
        day: &NaiveDate,
    ) -> Result<Vec<models::Booking>, Status> {
        Ok(self
            .repository
            .get_bookings(
                BookingsFilter {
                    venue: Some(Uuid::parse_str(venue_id).map_err(|e| {
                        log::error!("could not parse uuid : {}", e);
                        Status::internal("could not parse uuid")
                    })?),
                    day: Some(*day),
                },
                None,
                None,
            )?
            .to_vec())
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
            .find(|&hours| hours.day_of_week == date.weekday().number_from_monday())
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

#[cfg_attr(test, automock)]
pub trait UuidGetter {
    fn uuid(&self) -> Uuid;
}

#[derive(Default)]
struct GetUuid {}

impl UuidGetter for GetUuid {
    fn uuid(&self) -> Uuid {
        Uuid::new_v4()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::models::Booking;
    use chrono::{DateTime, Duration, NaiveDateTime, Utc};
    use mockall::predicate;
    use protobuf::venue::models::OpeningHoursSpecification;
    use std::convert::TryInto;
    use uuid::Uuid;

    #[test]
    fn test_get_free_table_no_bookings() {
        let free_table = super::get_free_table(
            60,
            &mut vec!["3a3789ca-7174-4127-ae50-a644d69f1d27".to_string()],
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
            &mut vec!["3a3789ca-7174-4127-ae50-a644d69f1d27".to_string()],
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
                name: "".to_string(),
                given_name: "".to_string()
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
            &mut vec!["3a3789ca-7174-4127-ae50-a644d69f1d27".to_string()],
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
                name: "".to_string(),
                given_name: "".to_string()
            }],
            &starts_at,
        );

        assert_eq!(free_table, None)
    }

    #[test]
    fn test_get_free_table_with_no_tables() {
        let starts_at = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(61, 0), Utc);
        let free_table = super::get_free_table(60, &mut vec![], &vec![], &starts_at);

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
                    slug: "".to_string(),
                })
            });
        let service = BookingService::new(Box::new(MockRepository::new()), Box::new(mock), None)
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

    #[tokio::test]
    async fn test_get_slot() {
        let starts = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);

        let venue_id = "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string();
        let people = 4;
        let duration = 60;

        let mut venue = MockVenueClient::new();
        let mut repository = MockRepository::new();

        venue
            .expect_get_venue()
            .with(predicate::eq(venue_id.clone()))
            .times(1)
            .returning(|_| {
                Ok(Venue {
                    id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                    name: "test venue".to_string(),
                    opening_hours: vec![OpeningHoursSpecification {
                        day_of_week: 5,
                        opens: "14:00".to_string(),
                        closes: "16:00".to_string(),
                        valid_from: "".to_string(),
                        valid_through: "".to_string(),
                    }],
                    special_opening_hours: vec![],
                    slug: "test-venue".to_string(),
                })
            });

        venue
            .expect_get_tables_with_capacity()
            .with(predicate::eq(venue_id.clone()), predicate::eq(people))
            .times(1)
            .returning(|_, _| Ok(vec!["eb7a8544-1595-4b62-ab72-137dd03b538f".to_string()]));

        repository
            .expect_get_bookings()
            .with(
                predicate::eq(BookingsFilter {
                    venue: Some(
                        Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                    ),
                    day: Some(starts.naive_utc().date()),
                }),
                predicate::eq(None),
                predicate::eq(None),
            )
            .times(1)
            .returning(|_, _, _| Ok(vec![]));

        let service = BookingService::new(Box::new(repository), Box::new(venue), None)
            .expect("could not construct booking service");

        let result = service
            .get_slot(Request::new(SlotInput {
                venue_id,
                email: "test@test.com".to_string(),
                people,
                starts_at: starts.to_rfc3339(),
                duration,
            }))
            .await
            .map(|r| r.into_inner())
            .expect("did not expect error from get slot");

        assert_eq!(
            result,
            GetSlotResponse {
                r#match: Some(Slot {
                    venue_id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                    email: "test@test.com".to_string(),
                    people: 4,
                    starts_at: "1992-05-01T15:00:00+00:00".to_string(),
                    ends_at: "1992-05-01T16:00:00+00:00".to_string(),
                    duration: 60
                }),
                other_available_slots: vec![
                    Slot {
                        venue_id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                        email: "test@test.com".to_string(),
                        people: 4,
                        starts_at: "1992-05-01T14:00:00+00:00".to_string(),
                        ends_at: "1992-05-01T15:00:00+00:00".to_string(),
                        duration: 60
                    },
                    Slot {
                        venue_id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                        email: "test@test.com".to_string(),
                        people: 4,
                        starts_at: "1992-05-01T14:30:00+00:00".to_string(),
                        ends_at: "1992-05-01T15:30:00+00:00".to_string(),
                        duration: 60
                    },
                    Slot {
                        venue_id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                        email: "test@test.com".to_string(),
                        people: 4,
                        starts_at: "1992-05-01T15:00:00+00:00".to_string(),
                        ends_at: "1992-05-01T16:00:00+00:00".to_string(),
                        duration: 60
                    },
                ]
            }
        )
    }

    #[tokio::test]
    async fn test_create_booking() {
        let starts = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);

        let venue_id = "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string();
        let people = 4;
        let duration = 60;

        let mut venue = MockVenueClient::new();
        let mut repository = MockRepository::new();

        venue
            .expect_get_venue()
            .with(predicate::eq(venue_id.clone()))
            .times(1)
            .returning(|_| {
                Ok(Venue {
                    id: "3a3789ca-7174-4127-ae50-a644d69f1d27".to_string(),
                    name: "test venue".to_string(),
                    opening_hours: vec![OpeningHoursSpecification {
                        day_of_week: 5,
                        opens: "14:00".to_string(),
                        closes: "16:00".to_string(),
                        valid_from: "".to_string(),
                        valid_through: "".to_string(),
                    }],
                    special_opening_hours: vec![],
                    slug: "test-venue".to_string(),
                })
            });

        venue
            .expect_get_tables_with_capacity()
            .with(predicate::eq(venue_id.clone()), predicate::eq(people))
            .times(1)
            .returning(|_, _| Ok(vec!["eb7a8544-1595-4b62-ab72-137dd03b538f".to_string()]));

        repository
            .expect_get_bookings()
            .with(
                predicate::eq(BookingsFilter {
                    venue: Some(
                        Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                    ),
                    day: Some(starts.naive_utc().date()),
                }),
                predicate::eq(None),
                predicate::eq(None),
            )
            .times(1)
            .returning(|_, _, _| Ok(vec![]));

        repository
            .expect_create_booking()
            .with(predicate::eq(Booking {
                id: Uuid::parse_str("5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d")
                    .expect("could not parse mock uuid"),
                customer_email: "test@test.com".to_string(),
                venue_id: Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                table_id: Uuid::parse_str("eb7a8544-1595-4b62-ab72-137dd03b538f")
                    .expect("could not parse table uuid"),
                people: people.try_into().unwrap(),
                date: starts.naive_utc().date(),
                starts_at: starts,
                ends_at: starts + Duration::minutes(duration),
                duration: duration as i32,
                name: "matthew cobbing".to_string(),
                given_name: "matthew".to_string()
            }))
            .times(1)
            .returning(|_| Ok(()));

        let mut get_uuid = MockUuidGetter::new();
        get_uuid.expect_uuid().returning(|| {
            Uuid::parse_str("5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d")
                .expect("could not parse mock uuid")
        });

        let service = BookingService::new(
            Box::new(repository),
            Box::new(venue),
            Some(Box::new(get_uuid)),
        )
        .expect("could not construct booking service");

        let result = service
            .create_booking(Request::new(BookingInput {
                venue_id: venue_id.clone(),
                email: "test@test.com".to_string(),
                people,
                starts_at: starts.to_rfc3339(),
                duration: duration as u32,
                name: "matthew cobbing".to_string(),
                given_name: "matthew".to_string()
            }))
            .await
            .map(|r| r.into_inner())
            .expect("did not expect error from create booking");

        assert_eq!(
            result,
            protobuf::booking::models::Booking {
                id: "5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d".to_string(),
                venue_id,
                email: "test@test.com".to_string(),
                people,
                starts_at: "1992-05-01T15:00:00+00:00".to_string(),
                ends_at: "1992-05-01T16:00:00+00:00".to_string(),
                duration: duration as u32,
                table_id: "eb7a8544-1595-4b62-ab72-137dd03b538f".to_string(),
                name: "matthew cobbing".to_string(),
                given_name: "matthew".to_string()
            }
        );
    }

    #[tokio::test]
    async fn test_get_bookings() {
        let starts = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);
        let mut mock = MockRepository::new();

        let venue_id = "0441d2c0-458d-4f4a-81cc-92e44807365b".to_string();

        mock.expect_count_bookings()
            .with(predicate::eq(BookingsFilter {
                venue: Some(
                    Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                ),
                day: Some(starts.naive_utc().date()),
            }))
            .times(1)
            .returning(|_| Ok(3));

        mock.expect_get_bookings()
            .with(
                predicate::eq(BookingsFilter {
                    venue: Some(
                        Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                    ),
                    day: Some(starts.naive_utc().date()),
                }),
                predicate::eq(Some(0)),
                predicate::eq(Some(2)),
            )
            .times(1)
            .returning(|filter, _, _| {
                let starts =
                    DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);
                let uuid = "5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d".to_string();
                Ok(vec![
                    models::Booking {
                        id: Uuid::parse_str(&uuid).unwrap(),
                        customer_email: "test@test.com".to_string(),
                        venue_id: filter.venue.unwrap(),
                        table_id: Uuid::parse_str(&uuid).unwrap(),
                        people: 4,
                        date: starts.naive_utc().date(),
                        starts_at: starts,
                        ends_at: starts + Duration::minutes(30),
                        duration: 30,
                        name: "matthew cobbing".to_string(),
                        given_name: "matthew".to_string()
                    },
                    models::Booking {
                        id: Uuid::parse_str(&uuid).unwrap(),
                        customer_email: "test2@test.com".to_string(),
                        venue_id: filter.venue.unwrap(),
                        table_id: Uuid::parse_str(&uuid).unwrap(),
                        people: 4,
                        date: starts.naive_utc().date(),
                        starts_at: starts,
                        ends_at: starts + Duration::minutes(30),
                        duration: 30,
                        name: "matthew cobbing".to_string(),
                        given_name: "matthew".to_string()
                    },
                    models::Booking {
                        id: Uuid::parse_str(&uuid).unwrap(),
                        customer_email: "test3@test.com".to_string(),
                        venue_id: filter.venue.unwrap(),
                        table_id: Uuid::parse_str(&uuid).unwrap(),
                        people: 4,
                        date: starts.naive_utc().date(),
                        starts_at: starts,
                        ends_at: starts + Duration::minutes(30),
                        duration: 30,
                        name: "matthew cobbing".to_string(),
                        given_name: "matthew".to_string()
                    },
                ])
            });

        let service = BookingService::new(Box::new(mock), Box::new(MockVenueClient::new()), None)
            .expect("could not construct booking service");

        let result = service
            .get_bookings(Request::new(GetBookingsRequest {
                venue_id: venue_id.clone(),
                date: starts.to_rfc3339(),
                page: 0,
                limit: 2,
            }))
            .await
            .map(|r| r.into_inner())
            .expect("did not expect error from get bookings");

        assert_eq!(
            result,
            GetBookingsResponse {
                bookings: vec![
                    protobuf::booking::models::Booking {
                        id: "5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d".to_string(),
                        venue_id: venue_id.clone(),
                        email: "test@test.com".to_string(),
                        people: 4,
                        starts_at: "1992-05-01T15:00:00+00:00".to_string(),
                        ends_at: "1992-05-01T15:30:00+00:00".to_string(),
                        duration: 30,
                        table_id: "5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d".to_string(),
                        name: "matthew cobbing".to_string(),
                        given_name: "matthew".to_string()
                    },
                    protobuf::booking::models::Booking {
                        id: "5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d".to_string(),
                        venue_id: venue_id.clone(),
                        email: "test2@test.com".to_string(),
                        people: 4,
                        starts_at: "1992-05-01T15:00:00+00:00".to_string(),
                        ends_at: "1992-05-01T15:30:00+00:00".to_string(),
                        duration: 30,
                        table_id: "5a77fdd3-9f2c-4096-8fc3-8eaae0d54e1d".to_string(),
                        name: "matthew cobbing".to_string(),
                        given_name: "matthew".to_string()
                    }
                ],
                has_next_page: true,
                pages: 2
            }
        )
    }

    #[tokio::test]
    async fn test_get_bookings_no_bookings() {
        let starts = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);
        let mut mock = MockRepository::new();

        let venue_id = "0441d2c0-458d-4f4a-81cc-92e44807365b".to_string();

        mock.expect_count_bookings()
            .with(predicate::eq(BookingsFilter {
                venue: Some(
                    Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                ),
                day: Some(starts.naive_utc().date()),
            }))
            .times(1)
            .returning(|_| Ok(0));

        mock.expect_get_bookings()
            .with(
                predicate::eq(BookingsFilter {
                    venue: Some(
                        Uuid::parse_str(&venue_id.clone()).expect("could not parse venue uuid"),
                    ),
                    day: Some(starts.naive_utc().date()),
                }),
                predicate::eq(Some(0)),
                predicate::eq(Some(5)),
            )
            .times(1)
            .returning(|_, _, _| Ok(vec![]));

        let service = BookingService::new(Box::new(mock), Box::new(MockVenueClient::new()), None)
            .expect("could not construct booking service");

        let result = service
            .get_bookings(Request::new(GetBookingsRequest {
                venue_id,
                date: starts.to_rfc3339(),
                page: 0,
                limit: 5,
            }))
            .await
            .map(|r| r.into_inner())
            .expect("did not expect error from get bookings");

        assert_eq!(
            result,
            GetBookingsResponse {
                bookings: vec![],
                has_next_page: false,
                pages: 0
            }
        )
    }

    #[tokio::test]
    async fn test_cancel_booking() {
        let starts = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);
        let mut mock = MockRepository::new();

        let venue_id = "0441d2c0-458d-4f4a-81cc-92e44807365b".to_string();
        let table_id = "9dce92ae-1605-4045-8ec5-ba79ad2165b0".to_string();
        let booking_id = "150dfddc-516f-443f-bb48-6add094e545d".to_string();

        mock.expect_cancel_booking()
            .with(
                predicate::eq(Uuid::parse_str(&booking_id.clone()).expect("could not parse venue uuid")),
            )
            .times(1)
            .returning(|_| {
                let starts = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(704732400, 0), Utc);
                Ok(models::Booking{
                    id: Uuid::parse_str("150dfddc-516f-443f-bb48-6add094e545d").expect("could not parse venue uuid"),
                    customer_email: "test@test.com".to_string(),
                    venue_id: Uuid::parse_str("0441d2c0-458d-4f4a-81cc-92e44807365b").expect("could not parse venue uuid"),
                    table_id: Uuid::parse_str("9dce92ae-1605-4045-8ec5-ba79ad2165b0").expect("could not parse venue uuid"),
                    people: 4,
                    date: starts.naive_utc().date(),
                    starts_at: starts,
                    ends_at: starts + Duration::minutes(60),
                    duration: 60,
                    name: "matthew cobbing".to_string(),
                    given_name: "matthew".to_string()
                })
            });

        let service = BookingService::new(Box::new(mock), Box::new(MockVenueClient::new()), None)
            .expect("could not construct booking service");

        let result = service
            .cancel_booking(Request::new(CancelBookingRequest{ id: booking_id.clone() }))
            .await
            .map(|r| r.into_inner())
            .expect("did not expect error from cancel booking");

        assert_eq!(
            result,
            protobuf::booking::models::Booking{
                id: booking_id,
                venue_id,
                email: "test@test.com".to_string(),
                people: 4,
                starts_at: starts.to_rfc3339(),
                ends_at: (starts + Duration::minutes(60)).to_rfc3339(),
                duration: 60,
                table_id,
                name: "matthew cobbing".to_string(),
                given_name: "matthew".to_string()
            }
        )
    }
}
