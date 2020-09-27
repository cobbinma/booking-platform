package services

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/models"
)

type createBookingService struct {
	repository  models.Repository
	tableClient models.TableClient
}

func NewCreateBookingService(repository models.Repository, tableClient models.TableClient) *createBookingService {
	return &createBookingService{repository: repository, tableClient: tableClient}
}

func (cbs *createBookingService) CreateBooking(ctx context.Context, slot models.NewBooking) (*models.Booking, error) {
	if err := slot.Valid(ctx, cbs.tableClient); err != nil {
		return nil, fmt.Errorf("slot is not valid : %s : %w", err, models.ErrInvalidRequest)
	}

	bookings, err := cbs.repository.GetBookings(ctx, models.BookingFilterWithTableIDs([]models.TableID{slot.TableID}), models.BookingFilterWithDate(&slot.Date))
	if err != nil {
		return nil, fmt.Errorf("could not get bookings : %s : %w", err, models.ErrInternalError)
	}

	for i := range bookings {
		if slot.StartsAt.Before(bookings[i].EndsAt) && bookings[i].StartsAt.Before(slot.EndsAt) {
			err := "incorrect user request : requested booking slot is not free"
			return nil, fmt.Errorf("%s , %w", err, models.ErrInvalidRequest)
		}
	}

	booking, err := cbs.repository.CreateBooking(ctx, slot)
	if err != nil {
		return nil, fmt.Errorf("repository error : %s : %w", err, models.ErrInternalError)
	}

	return booking, nil
}
