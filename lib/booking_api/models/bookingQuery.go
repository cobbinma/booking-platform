package models

import "time"

type BookingQuery struct {
	People   int       `json:"people" db:"people"`
	Date     time.Time `json:"date" db:"date"`
	StartsAt time.Time `json:"starts_at" db:"starts_at"`
	EndsAt   time.Time `json:"ends_at" db:"ends_at"`
}
