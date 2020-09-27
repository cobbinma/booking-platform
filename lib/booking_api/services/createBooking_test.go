package services_test

import (
	"context"
	"github.com/cobbinma/booking/lib/booking_api/models"
	mockmodels "github.com/cobbinma/booking/lib/booking_api/models/mock"
	"github.com/cobbinma/booking/lib/booking_api/services"
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
		repository  *mockmodels.MockRepository
		tableClient *mockmodels.MockTableClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		repository = mockmodels.NewMockRepository(ctrl)
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
	})
})
