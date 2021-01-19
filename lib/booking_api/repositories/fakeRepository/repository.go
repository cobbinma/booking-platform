package fakeRepository

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/booking_api/models"
)

type fakeRepository struct {
	bookings map[int]*models.Booking
}

func NewFakeRepository() models.Repository {
	return &fakeRepository{
		bookings: make(map[int]*models.Booking),
	}
}

func (f *fakeRepository) Migrate(ctx context.Context, sourceURL string) error {
	panic("implement me")
}

func (f *fakeRepository) CreateBooking(ctx context.Context, booking models.Slot) (*models.Booking, error) {
	newBooking := &models.Booking{
		ID:         len(f.bookings),
		CustomerID: booking.CustomerID,
		TableID:    booking.TableID,
		People:     booking.People,
		Date:       booking.Date,
		StartsAt:   booking.StartsAt,
		EndsAt:     booking.EndsAt,
	}
	f.bookings[len(f.bookings)] = newBooking
	return newBooking, nil
}

func (f *fakeRepository) GetBookings(ctx context.Context, options ...func(*models.BookingFilter) *models.BookingFilter) ([]models.Booking, error) {
	bookings := []models.Booking{}

	filter := &models.BookingFilter{}
	for _, option := range options {
		option(filter)
	}

	for _, booking := range f.bookings {
		if filter.Date != nil && *filter.Date != booking.Date {
			continue
		}

		add := true
		if len(filter.TableIDs) > 0 {
			for i := range filter.TableIDs {
				if booking.TableID == filter.TableIDs[i] {
					break
				}
				add = false
			}
		}

		if add {
			bookings = append(bookings, *booking)
		}
	}

	return bookings, nil
}

func (f *fakeRepository) DeleteBookings(ctx context.Context, id []int) error {
	for i := range id {
		delete(f.bookings, id[i])
	}

	return nil
}
