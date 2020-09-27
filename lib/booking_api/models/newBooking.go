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

	if err := dateTimesValidator(nb.Date, nb.StartsAt, nb.EndsAt); err != nil {
		return err
	}

	venue, ok := ctx.Value(VenueCtxKey).(Venue)
	if !ok {
		return fmt.Errorf("venue was not in context")
	}
	if !venue.IsOpen(nb.Date.Time().Day(), nb.StartsAt, nb.EndsAt) {
		return fmt.Errorf("venue is not open at those times")
	}

	table, err := tc.GetTable(ctx, nb.TableID)
	if err != nil {
		return fmt.Errorf("could not find table : %w", err)
	}

	if !table.HasCapacity(nb.People) {
		return fmt.Errorf("requested table does not have capacity")
	}

	return nil
}

func dateTimesValidator(date Date, startsAt time.Time, endsAt time.Time) error {
	now := time.Now()
	if date.Time().Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())) {
		return fmt.Errorf("date must not be in the past")
	}

	if !sameDate(date.Time(), startsAt) || !sameDate(date.Time(), endsAt) || !sameDate(endsAt, startsAt) {
		return fmt.Errorf("all times must be on same date")
	}

	if startsAt.After(endsAt) {
		return fmt.Errorf("starts at cannot be after ends at")
	}

	if startsAt.Before(time.Date(
		date.Time().Year(), date.Time().Month(), date.Time().Day(), 0, 0, 0, -1, startsAt.Location())) {
		return fmt.Errorf("starts at must be after date")
	}

	if endsAt.After(startsAt.Add(12*time.Hour + time.Second)) {
		return fmt.Errorf("booking can not exceed 12 hours")
	}

	return nil
}

func sameDate(a time.Time, b time.Time) bool {
	if a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year() {
		return true
	}
	return false
}
