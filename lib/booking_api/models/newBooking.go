package models

import (
	"context"
	"fmt"
	"time"
)

type NewBooking struct {
	CustomerID CustomerID `json:"customer_id" db:"customer_id"`
	TableID    TableID    `json:"table_id" db:"table_id"`
	People     int        `json:"people" db:"people"`
	Date       Date       `json:"date" db:"date"`
	StartsAt   time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt     time.Time  `json:"ends_at" db:"ends_at"`
}

func (nb NewBooking) Valid(ctx context.Context, tc TableClient) error {
	if err := nb.CustomerID.Valid(); err != nil {
		return err
	}

	if err := nb.TableID.Valid(); err != nil {
		return err
	}

	if nb.People < 1 {
		return fmt.Errorf("must have positive people")
	}

	if err := nb.dateTimesValidator(); err != nil {
		return err
	}

	table, err := tc.GetTable(ctx, nb.TableID)
	if err != nil {
		return fmt.Errorf("could not find table : %w", err)
	}

	if table.Capacity < nb.People {
		return fmt.Errorf("requested table does not have capacity")
	}

	return nil
}

func (nb NewBooking) dateTimesValidator() error {
	now := time.Now()
	if nb.Date.Time().Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())) {
		return fmt.Errorf("date must not be in the past")
	}

	if nb.StartsAt.After(nb.EndsAt) {
		return fmt.Errorf("starts at cannot be after ends at")
	}

	if nb.StartsAt.Before(time.Date(nb.Date.Time().Year(), nb.Date.Time().Month(), nb.Date.Time().Day(), 0, 0, 0, -1, nb.StartsAt.Location())) {
		return fmt.Errorf("starts at must be after date")
	}

	if nb.EndsAt.After(nb.StartsAt.Add(12*time.Hour + time.Second)) {
		return fmt.Errorf("booking can not exceed 12 hours")
	}

	return nil
}
