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
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return fmt.Errorf("venue was not in context")
	}

	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("bookings").
		Columns("venue_id", "customer_id", "table_id", "people", "date", "starts_at", "ends_at").
		Values(venue.ID, booking.CustomerID, booking.TableID, booking.People, booking.Date.Time(), booking.StartsAt, booking.EndsAt).
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

func (p *postgres) GetBookings(ctx context.Context, options ...func(*models.BookingFilter) *models.BookingFilter) ([]models.Booking, error) {
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return nil, fmt.Errorf("venue was not in context")
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "customer_id", "table_id", "people", "date", "starts_at", "ends_at").
		From("bookings")

	filter := &models.BookingFilter{}
	for _, option := range options {
		option(filter)
	}

	where := sq.And{}
	if filter.Date != nil {
		where = append(where, sq.Eq{"date": filter.Date})
	}
	if filter.TableIDs != nil {
		where = append(where, sq.Eq{"table_id": filter.TableIDs})
	}

	where = append(where, sq.Eq{"venue_id": venue.ID})

	builder = builder.Where(where)

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

func (p *postgres) DeleteBookings(ctx context.Context, id []int) error {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("bookings").
		Where(sq.Eq{"id": id}).ToSql()
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
