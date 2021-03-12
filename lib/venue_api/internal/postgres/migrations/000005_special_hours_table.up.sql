CREATE TABLE IF NOT EXISTS special_opening_hours
(
    venue_id       UUID NOT NULL REFERENCES venues (id) ON DELETE CASCADE,
    day_of_week    INTEGER NOT NULL,
    opens          VARCHAR NOT NULL,
    closes         VARCHAR NOT NULL,
    valid_from     DATE NOT NULL,
    valid_through  DATE NOT NULL,
    CHECK (day_of_week > 0),
    CHECK (day_of_week < 8),
    CHECK (valid_through > valid_from)
);