CREATE TABLE IF NOT EXISTS venues (
     id  SERIAL UNIQUE PRIMARY KEY,
     day_of_week INTEGER NOT NULL,
     opens TIME NOT NULL,
     closes TIME NOT NULL
);

CREATE TABLE IF NOT EXISTS opening_hours (
     id  SERIAL UNIQUE PRIMARY KEY,
     venue_id   INTEGER,
     name VARCHAR NOT NULL,
     capacity INTEGER NOT NULL,
     CONSTRAINT fk_openhours_venue
         FOREIGN KEY (venue_id)
             REFERENCES venues(id) ON DELETE CASCADE
);