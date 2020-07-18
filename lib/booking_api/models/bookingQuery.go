package models

type BookingQuery struct {
	People   int    `json:"people" db:"people"`
	StartsAt string `json:"starts_at" db:"starts_at"`
	EndsAt   string `json:"ends_at" db:"ends_at"`
}
