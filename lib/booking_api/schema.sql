CREATE TABLE IF NOT EXISTS bookings (
    id  uuid PRIMARY KEY NOT NULL,
    customer_email VARCHAR NOT NULL,
    venue_id uuid NOT NULL,
    table_id uuid NOT NULL,
    people  INTEGER NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration INTEGER NOT NULL,
    CHECK (people > 0),
    CHECK (ends_at > starts_at),
    CHECK (starts_at > NOW()),
    CHECK (ends_at > NOW()),
    CHECK (duration > 0)
);