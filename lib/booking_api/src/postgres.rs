#[warn(unused_imports)]
use crate::diesel::RunQueryDsl;

use crate::models::Booking;
use chrono::{DateTime, Datelike, Duration, NaiveDate};
use diesel::{Connection, ExpressionMethods, PgConnection, QueryDsl};
use protobuf::booking::models::SlotInput;
use std::env;
use tonic::{Request, Status};
use uuid::Uuid;

struct Postgres {
    connection: PgConnection,
}

impl Postgres {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let connection = PgConnection::establish(&database_url)
            .expect(&format!("Error connecting to {}", database_url));

        Ok(Postgres { connection })
    }

    pub fn get_bookings(&self, req: Request<SlotInput>) -> Result<Vec<Booking>, Status> {
        use crate::schema::bookings::dsl::*;
        let slot = req.into_inner();
        let s = DateTime::parse_from_rfc3339(&slot.starts_at)
            .map_err(|_| Status::internal("could not parse date"))?;
        let day = NaiveDate::from_ymd(s.year(), s.month(), s.day());

        let results: Vec<Booking> = bookings
            .filter(
                venue_id.eq(Uuid::parse_str(&slot.venue_id)
                    .map_err(|_| Status::internal("could not parse uuid"))?),
            )
            .filter(date.gt(day))
            .filter(date.lt(day + Duration::days(1)))
            .get_results(&self.connection)
            .map_err(|_| Status::internal("could not get get bookings"))?;

        Ok(results)
    }
}
