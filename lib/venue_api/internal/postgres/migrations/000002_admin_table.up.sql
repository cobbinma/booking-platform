CREATE TABLE IF NOT EXISTS admins
(
    id   UUID UNIQUE PRIMARY KEY NOT NULL,
    venue_id UUID NOT NULL,
    email VARCHAR NOT NULL,
    UNIQUE(venue_id, email)
);