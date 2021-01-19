package services_test

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/booking_api/models"
	mockmodels "github.com/cobbinma/booking-platform/lib/booking_api/models/mock"
	"github.com/cobbinma/booking-platform/lib/booking_api/repositories/fakeRepository"
	"github.com/cobbinma/booking-platform/lib/booking_api/services"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var (
	now         = time.Now()
	twoDays     = time.Date(now.Year(), now.Month(), now.Day()+2, 0, 0, 0, 0, time.UTC)
	twoDaysDate = models.Date(twoDays)
	startsAt    = time.Date(now.Year(), now.Month(), now.Day()+2, 18, 0, 0, 0, time.UTC)
	endsAt      = time.Date(now.Year(), now.Month(), now.Day()+2, 20, 0, 0, 0, time.UTC)
)

var _ = Describe("CreateCoupon", func() {
	var (
		ctx         context.Context
		ctrl        *gomock.Controller
		repository  models.Repository
		tableClient *mockmodels.MockTableClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		repository = fakeRepository.NewFakeRepository()
		tableClient = mockmodels.NewMockTableClient(ctrl)
	})

	Context("with valid slot", func() {
	})

	Context("with invalid slot", func() {
		When("venue is shut", func() {
			slot := models.Slot{
				CustomerID: "test@test.test",
				TableID:    1,
				People:     4,
				Date:       twoDaysDate,
				StartsAt:   startsAt,
				EndsAt:     endsAt,
			}
			venue := models.Venue{
				ID:           1,
				Name:         "Hop and Vine",
				OpeningHours: []models.OpeningHours{},
			}
			BeforeEach(func() {
				ctx = context.WithValue(ctx, models.VenueCtxKey, venue)
			})
			It("should return a invalid request error", func() {
				service := services.NewCreateBookingService(repository, tableClient)
				_, err := service.CreateBooking(ctx, slot)
				Expect(err).Should(MatchError(models.ErrInvalidRequest))
			})
		})
		When("table does not have capacity", func() {
			tableID := models.TableID(1)
			table := models.Table{
				Name:     "small table",
				Capacity: 2,
			}
			slot := models.Slot{
				CustomerID: "test@test.test",
				TableID:    tableID,
				People:     4,
				Date:       twoDaysDate,
				StartsAt:   startsAt,
				EndsAt:     endsAt,
			}
			venue := models.Venue{
				ID:   1,
				Name: "Hop and Vine",
				OpeningHours: []models.OpeningHours{{
					DayOfWeek: twoDays.Day(),
					Opens:     models.TimeOfDay(startsAt),
					Closes:    models.TimeOfDay(endsAt),
				}},
			}
			BeforeEach(func() {
				ctx = context.WithValue(ctx, models.VenueCtxKey, venue)
				tableClient.EXPECT().GetTable(gomock.Eq(ctx), gomock.Eq(tableID)).Return(&table, nil)
			})
			It("should return a invalid request error", func() {
				service := services.NewCreateBookingService(repository, tableClient)
				_, err := service.CreateBooking(ctx, slot)
				Expect(err).Should(MatchError(models.ErrInvalidRequest))
			})
		})
	})
})
