use super::schema::bookings;
use chrono::{DateTime, NaiveDate, Utc};
use uuid::Uuid;

#[derive(Queryable, Debug, Insertable, Clone)]
#[table_name = "bookings"]
pub struct Booking {
    pub id: Uuid,
    pub customer_email: String,
    pub venue_id: Uuid,
    pub table_id: Uuid,
    pub people: i32,
    pub date: NaiveDate,
    pub starts_at: DateTime<Utc>,
    pub ends_at: DateTime<Utc>,
    pub duration: i32,
}
