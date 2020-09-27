package models

import (
	"context"
	"fmt"
	"time"
)

type BookingQuery struct {
	CustomerID CustomerID `json:"customer_id" db:"customer_id"`
	People     int        `json:"people" db:"people"`
	Date       Date       `json:"date" db:"date"`
	StartsAt   time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt     time.Time  `json:"ends_at" db:"ends_at"`
}

func (bq BookingQuery) Valid(ctx context.Context) error {
	if err := bq.CustomerID.Valid(); err != nil {
		return err
	}

	if bq.People < 1 {
		return fmt.Errorf("must have positive people")
	}

	if err := dateTimesValidator(bq.Date, bq.StartsAt, bq.EndsAt); err != nil {
		return err
	}

	venue, ok := ctx.Value(VenueCtxKey).(Venue)
	if !ok {
		return fmt.Errorf("venue was not in context")
	}
	if !venue.IsOpen(bq.Date.Time().Day(), bq.StartsAt, bq.EndsAt) {
		return fmt.Errorf("venue is not open at those times")
	}

	return nil
}
