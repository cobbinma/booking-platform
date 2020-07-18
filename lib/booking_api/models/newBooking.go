package models

import (
	"fmt"
	"time"
)

type NewBooking struct {
	CustomerID CustomerID `json:"customer_id" db:"customer_id"`
	TableID    TableID    `json:"table_id" db:"table_id"`
	People     int        `json:"people" db:"people"`
	Date       time.Time  `json:"date" db:"date"`
	StartsAt   time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt     time.Time  `json:"ends_at" db:"ends_at"`
}

func (nb NewBooking) Valid() error {
	if err := nb.CustomerID.Valid(); err != nil {
		return err
	}

	if err := nb.TableID.Valid(); err != nil {
		return err
	}

	if nb.People < 1 {
		return fmt.Errorf("must have positive people")
	}

	now := time.Now()
	if nb.Date.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)) {
		return fmt.Errorf("date must not be in the past")
	}

	if nb.StartsAt.After(nb.EndsAt) {
		return fmt.Errorf("starts at cannot be after ends at")
	}

	if nb.StartsAt.Before(time.Date(nb.Date.Year(), nb.Date.Month(), nb.Date.Day(), 0, 0, 0, 0, time.UTC)) {
		return fmt.Errorf("starts at must be after date")
	}

	if nb.EndsAt.After(nb.StartsAt.Add(12*time.Hour + time.Second)) {
		return fmt.Errorf("booking can not be longer than 12 hours")
	}

	return nil
}
