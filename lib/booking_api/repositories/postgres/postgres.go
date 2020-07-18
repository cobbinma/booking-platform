package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/cobbinma/booking/lib/booking_api/models"
)

type postgres struct {
	dbClient DBClient
}

func (p *postgres) CreateBooking(ctx context.Context, booking models.NewBooking) error {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("bookings").
		Columns("customer_id", "table_id", "people", "date", "starts_at", "ends_at").
		Values(booking.CustomerID, booking.TableID, booking.People, booking.Date.Time(), booking.StartsAt, booking.EndsAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s : %w", "could not build statement", err)
	}

	_, err = p.dbClient.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("%s : %w", "could not execute", err)
	}

	return nil
}

func NewPostgres(client DBClient) models.Repository {
	return &postgres{dbClient: client}
}
