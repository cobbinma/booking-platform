CREATE TABLE IF NOT EXISTS venues (
    id  SERIAL UNIQUE PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS opening_hours (
    id  SERIAL UNIQUE PRIMARY KEY,
    venue_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL,
    opens TIME NOT NULL,
    closes TIME NOT NULL,
    UNIQUE (venue_id, day_of_week),
    CONSTRAINT fk_openhours_venue
     FOREIGN KEY (venue_id)
         REFERENCES venues(id) ON DELETE CASCADE
);