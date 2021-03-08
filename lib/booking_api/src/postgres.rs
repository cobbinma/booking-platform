#[warn(unused_imports)]
use crate::diesel::RunQueryDsl;

use crate::models::Booking;
use crate::service::{BookingsFilter, Repository};
use diesel::r2d2::{ConnectionManager, Pool};
use diesel::{ExpressionMethods, PgConnection, QueryDsl};
use std::env;
use tonic::Status;
use uuid::Uuid;

embed_migrations!("./migrations");

pub struct Postgres {
    pool: Pool<ConnectionManager<PgConnection>>,
}

impl Postgres {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let database_url = env::var("DATABASE_URL")?;
        let manager = ConnectionManager::<PgConnection>::new(database_url);
        let pool = diesel::r2d2::Builder::new().build(manager)?;

        embedded_migrations::run(&pool.get()?)?;

        Ok(Postgres { pool })
    }
}

impl Repository for Postgres {
    fn get_bookings(
        &self,
        filter: BookingsFilter,
        page: Option<i32>,
        limit: Option<i32>,
    ) -> Result<Vec<Booking>, Status> {
        tracing::debug!("get bookings from postgres");
        use crate::schema::bookings::dsl::*;

        let mut builder = bookings.into_boxed();

        if let Some(uuid) = filter.venue {
            builder = builder.filter(venue_id.eq(uuid))
        };

        if let Some(day) = filter.day {
            builder = builder.filter(date.eq(day))
        };

        if let (Some(page), Some(limit)) = (page, limit) {
            builder = builder.limit(limit as i64 + 1);
            builder = builder.offset((limit * page) as i64);
        };

        builder
            .order(starts_at.asc())
            .get_results(&self.pool.get().map_err(|e| {
                log::error!("could not get database connection : {}", e);
                Status::internal("could not get database connection")
            })?)
            .map_err(|e| {
                log::error!("could not get get bookings from database : {}", e);
                Status::internal("could not get get bookings from database")
            })
    }

    fn create_booking(&self, new_booking: &Booking) -> Result<(), Status> {
        tracing::debug!(
            "create booking '{}' in postgres for venue '{}' on date '{}', starting at '{}', for '{}' minutes, for '{}' people",
            &new_booking.id.to_string(),
            &new_booking.venue_id.to_string(),
            &new_booking.date.to_string(),
            &new_booking.starts_at.to_rfc3339(),
            &new_booking.duration,
            &new_booking.people,
        );
        use crate::schema::bookings;
        diesel::insert_into(bookings::table)
            .values(new_booking)
            .execute(&self.pool.get().map_err(|e| {
                log::error!("{}", e);
                Status::internal("could not get database connection")
            })?)
            .map_err(|e| {
                log::error!("{}", e);
                Status::internal("could not create booking in database")
            })?;

        Ok(())
    }

    fn cancel_booking(&self, booking_id: &Uuid) -> Result<Booking, Status> {
        tracing::debug!(
            "cancelling booking '{}' from postgres",
            &booking_id.to_string()
        );
        use crate::schema::bookings::dsl::*;

        let booking = bookings
            .find(booking_id)
            .first(&self.pool.get().map_err(|e| {
                log::error!("could not get database connection : {}", e);
                Status::internal("could not get database connection")
            })?)
            .map_err(|e| {
                log::error!("could not get get booking from database : {}", e);
                Status::internal("could not get get booking from database")
            })?;

        diesel::delete(bookings.filter(id.eq(booking_id))).execute(&self.pool.get().map_err(|e| {
            log::error!("could not get database connection : {}", e);
            Status::internal("could not get database connection")
        })?)
            .map_err(|e| {
                log::error!("could not get delete booking from database : {}", e);
                Status::internal("could not delete booking from database")
            })?;

        Ok(booking)
    }

    fn count_bookings(&self, filter: &BookingsFilter) -> Result<i64, Status> {
        tracing::debug!("counting bookings from postgres");
        use crate::schema::bookings::dsl::*;

        let mut builder = bookings.into_boxed();

        if let Some(uuid) = filter.venue {
            builder = builder.filter(venue_id.eq(uuid))
        };

        if let Some(day) = filter.day {
            builder = builder.filter(date.eq(day))
        };

        builder
            .count()
            .get_result(&self.pool.get().map_err(|e| {
                log::error!("could not get database connection : {}", e);
                Status::internal("could not get database connection")
            })?)
            .map_err(|e| {
                log::error!("could not get get bookings from database : {}", e);
                Status::internal("could not get get bookings from database")
            })
    }
}
