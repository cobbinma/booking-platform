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

func (p *postgres) GetBookings(ctx context.Context, filter *models.BookingFilter) ([]models.Booking, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "customer_id", "table_id", "people", "date", "starts_at", "ends_at").
		From("bookings")

	if filter != nil && filter.Date != nil {
		builder = builder.Where(sq.Eq{"date": filter.Date})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not build statement", err)
	}

	tables, err := p.dbClient.GetBookings(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not get tables from db client", err)
	}

	return tables, nil
}

func NewPostgres(client DBClient) models.Repository {
	return &postgres{dbClient: client}
}
