package models

import "time"

type Booking struct {
	ID         int       `json:"id" db:"id"`
	CustomerID string    `json:"customer_id" db:"customer_id"`
	TableID    int       `json:"table_id" db:"table_id"`
	People     int       `json:"people" db:"people"`
	Date       time.Time `json:"date" db:"date"`
	StartsAt   time.Time `json:"starts_at" db:"starts_at"`
	EndsAt     time.Time `json:"ends_at" db:"ends_at"`
}

type NewBooking struct {
	CustomerID string    `json:"customer_id" db:"customer_id"`
	TableID    int       `json:"table_id" db:"table_id"`
	People     int       `json:"people" db:"people"`
	Date       time.Time `json:"date" db:"date"`
	StartsAt   time.Time `json:"starts_at" db:"starts_at"`
	EndsAt     time.Time `json:"ends_at" db:"ends_at"`
}
