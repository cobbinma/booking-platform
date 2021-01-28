CREATE TABLE IF NOT EXISTS venues
(
    id   UUID UNIQUE PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS tables
(
    id       UUID UNIQUE PRIMARY KEY NOT NULL REFERENCES venues (id) ON DELETE CASCADE,
    venue_id UUID NOT NULL,
    name     VARCHAR NOT NULL,
    capacity INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS opening_hours
(
    venue_id    UUID NOT NULL REFERENCES venues (id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL,
    opens       VARCHAR NOT NULL,
    closes      VARCHAR NOT NULL,
    UNIQUE (venue_id, day_of_week),
    CHECK (day_of_week > 0),
    CHECK (day_of_week < 8)
);