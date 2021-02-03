table! {
    bookings (id) {
        id -> Uuid,
        customer_email -> Varchar,
        venue_id -> Uuid,
        table_id -> Uuid,
        people -> Int4,
        date -> Date,
        starts_at -> Timestamptz,
        ends_at -> Timestamptz,
        duration -> Int4,
    }
}
