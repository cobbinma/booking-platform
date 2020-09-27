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

func (p *postgres) CreateBooking(ctx context.Context, booking models.Slot) (*models.Booking, error) {
	venue, ok := ctx.Value(models.VenueCtxKey).(models.Venue)
	if !ok {
		return nil, fmt.Errorf("venue was not in context")
	}

	args := map[string]interface{}{"venue_id": venue.ID, "customer_id": booking.CustomerID, "table_id": booking.TableID,
		"people": booking.People, "date": booking.Date.Time(), "starts_at": booking.StartsAt, "ends_at": booking.EndsAt}
	rows, err := p.dbClient.NamedQuery("INSERT INTO bookings (venue_id, customer_id, table_id, people, date, starts_at, ends_at) "+
		"VALUES (:venue_id, :customer_id, :table_id, :people, :date, :starts_at, :ends_at) RETURNING id", args)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "could not perform named query to insert into tables", err)
	}

	var bookingID int
	if rows.Next() {
		if err := rows.Scan(&bookingID); err != nil {
			return nil, fmt.Errorf("%s : %w", "could not scan row", err)
		}
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("%s : %w", "could not close rows", err)
	}

	return &models.Booking{
		ID:         bookingID,
		CustomerID: booking.CustomerID,
		TableID:    booking.TableID,
		People:     booking.People,
		Date:       booking.Date,
		StartsAt:   booking.StartsAt,
		EndsAt:     booking.EndsAt,
	}, nil
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
