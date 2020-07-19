CREATE TABLE IF NOT EXISTS bookings (
    id  SERIAL UNIQUE PRIMARY KEY,
    customer_id VARCHAR NOT NULL,
    table_id INTEGER NOT NULL,
    people  INTEGER NOT NULL,
    date    DATE    NOT NULL,
    starts_at TIMESTAMP NOT NULL,
    ends_at TIMESTAMP NOT NULL
);