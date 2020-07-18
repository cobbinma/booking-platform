package models

type Booking struct {
	ID         int    `json:"id" db:"id"`
	CustomerID string `json:"customer_id" db:"customer_id"`
	TableID    int    `json:"table_id" db:"table_id"`
	People     int    `json:"people" db:"people"`
	StartsAt   string `json:"starts_at" db:"starts_at"`
	EndsAt     string `json:"ends_at" db:"ends_at"`
}
