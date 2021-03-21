package graph_test

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/cmd/api/middleware"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	mock_resolver "github.com/cobbinma/booking-platform/lib/gateway_api/graph/mock"
	booking2 "github.com/cobbinma/booking-platform/lib/gateway_api/internal/booking"
	venue2 "github.com/cobbinma/booking-platform/lib/gateway_api/internal/venue"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	api2 "github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/booking/api"
	booking "github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/booking/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	venue "github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

//go:generate mockgen -package=mock_resolver -destination=./mock/venue.go github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api VenueAPIClient
//go:generate mockgen -package=mock_resolver -destination=./mock/booking.go github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/booking/api BookingAPIClient

func Test_GetVenue(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   venueID,
		Slug: "",
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Slug string `json:"slug"`
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(filter:{id:"a3291740-e89f-4cc0-845c-75c4c39842c9"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},slug}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetVenueTables(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	slug := "test-venue"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   "",
		Slug: slug,
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)

	venueClient.EXPECT().GetTables(gomock.Any(), &api.GetTablesRequest{VenueId: venueID}).Return(&api.GetTablesResponse{Tables: []*venue.Table{
		{
			Id:       "175fd06d-9a60-4ea6-86ca-bb96ca861208",
			Name:     "table one",
			Capacity: 4,
		},
	}}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Tables []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Capacity int    `json:"capacity"`
			} `json:"tables"`
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(filter:{slug:"test-venue"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},tables{id,name,capacity}}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetVenueTablesNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   venueID,
		Slug: "",
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{VenueId: venueID, Email: "test@test.com"}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Tables []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Capacity int    `json:"capacity"`
			} `json:"tables"`
		} `json:"getVenue"`
	}

	assert.Error(t, c.Post(`{getVenue(filter:{id:"a3291740-e89f-4cc0-845c-75c4c39842c9"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},tables{id,name,capacity}}}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetOpeningHoursSpecification(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	slug := "test-venue"
	date := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   "",
		Slug: slug,
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().GetOpeningHoursSpecification(gomock.Any(), &api.GetOpeningHoursSpecificationRequest{
		VenueId: venueID,
		Date:    date.Format(time.RFC3339),
	}).Return(&api.GetOpeningHoursSpecificationResponse{Specification: &venue.OpeningHoursSpecification{
		DayOfWeek:    1,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    "",
		ValidThrough: "",
	}}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			OpeningHoursSpecification struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHoursSpecification"`
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(filter:{slug:"test-venue"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},openingHoursSpecification(date:"3000-01-01T00:00:00Z"){dayOfWeek, opens, closes, validFrom, validThrough}}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetVenueAdmins(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	slug := "test-venue"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   "",
		Slug: slug,
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)

	venueClient.EXPECT().GetAdmins(gomock.Any(), &api.GetAdminsRequest{VenueId: venueID}).Return(&api.GetAdminsResponse{Admins: []string{"test@test.com"}}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Admins []string `json:"admins"`
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(filter:{slug:"test-venue"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},admins}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetVenueAdminsNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   venueID,
		Slug: "",
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Admins []string `json:"admins"`
		} `json:"getVenue"`
	}

	assert.Error(t, c.Post(`{getVenue(filter:{id:"a3291740-e89f-4cc0-845c-75c4c39842c9"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},admins}}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetVenueBookings(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	slug := "test-venue"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)
	bookingClient := mock_resolver.NewMockBookingAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   "",
		Slug: slug,
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)

	limit := 5
	date := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	bookingClient.EXPECT().GetBookings(gomock.Any(), &api2.GetBookingsRequest{
		VenueId: venueID,
		Date:    date.Format(time.RFC3339),
		Page:    0,
		Limit:   int32(limit),
	}).Return(&api2.GetBookingsResponse{
		Bookings: []*booking.Booking{
			{
				Id:        "cca3c988-9e11-4b81-9a98-c960fb4a3d97",
				VenueId:   "8a18e89b-339b-4e51-ab53-825aae59a070",
				Email:     "test@test.com",
				People:    5,
				StartsAt:  date.Format(time.RFC3339),
				EndsAt:    date.Add(time.Minute * 60).Format(time.RFC3339),
				Duration:  60,
				TableId:   "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
				Name:      "Test Test",
				GivenName: "Test",
			},
		},
		HasNextPage: false,
		Pages:       1,
	}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)
	bookingSrv, _, err := booking2.NewBookingClient("", nil, nil, booking2.WithClient(bookingClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, bookingSrv)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Bookings struct {
				Bookings []struct {
					ID        string `json:"id"`
					VenueID   string `json:"venueId"`
					Email     string `json:"email"`
					People    int    `json:"people"`
					StartsAt  string `json:"startsAt"`
					EndsAt    string `json:"endsAt"`
					Duration  int    `json:"duration"`
					TableID   string `json:"tableId"`
					Name      string `json:"name"`
					GivenName string `json:"givenName"`
				} `json:"bookings"`
				HasNextPage bool `json:"hasNextPage"`
				Pages       int  `json:"pages"`
			} `json:"bookings"`
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(filter:{slug:"test-venue"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},bookings(filter:{date:"3000-01-01T00:00:00Z"},pageInfo:{page:0,limit:5}){bookings{id,venueId,email,people,startsAt,endsAt,duration,tableId,name,givenName},hasNextPage,pages}}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_GetVenueBookingsNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().GetVenue(gomock.Any(), &api.GetVenueRequest{
		Id:   "",
		Slug: "test-venue",
	}).Return(&venue.Venue{
		Id:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        defaultOpeningHours(),
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetVenue struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				Opens        string `json:"opens"`
				Closes       string `json:"closes"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
			Bookings struct {
				Bookings []struct {
					ID       string `json:"id"`
					VenueID  string `json:"venueId"`
					Email    string `json:"email"`
					People   int    `json:"people"`
					StartsAt string `json:"startsAt"`
					EndsAt   string `json:"endsAt"`
					Duration int    `json:"duration"`
					TableID  string `json:"tableId"`
				} `json:"bookings"`
				HasNextPage bool `json:"hasNextPage"`
				Pages       int  `json:"pages"`
			} `json:"bookings"`
		} `json:"getVenue"`
	}

	assert.Error(t, c.Post(`{getVenue(filter:{slug:"test-venue"}){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough},bookings(filter:{date:"3000-01-01T00:00:00Z"},pageInfo:{page:0,limit:5}){bookings{id,venueId,email,people,startsAt,endsAt,duration,tableId},hasNextPage,pages}}}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_AddTableNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		AddTable struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Capacity int    `json:"capacity"`
		} `json:"addTable"`
	}
	assert.Error(t, c.Post(`mutation{addTable(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",name:"test table",capacity:5}) {id,name,capacity}}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_AddTable(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)
	venueClient.EXPECT().AddTable(gomock.Any(), &api.AddTableRequest{
		VenueId:  venueID,
		Name:     "test table",
		Capacity: 5,
	}).Return(&venue.Table{
		Id:       "bfcc0d78-83e7-4830-96ab-96cdbd0357c7",
		Name:     "test table",
		Capacity: 5,
	}, nil)

	var resp struct {
		AddTable struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Capacity int    `json:"capacity"`
		} `json:"addTable"`
	}

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	client.New(e).MustPost(`mutation{addTable(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",name:"test table",capacity:5}) {id,name,capacity}}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_UpdateOpeningHours(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)

	venueClient.EXPECT().UpdateOpeningHours(gomock.Any(), &api.UpdateOpeningHoursRequest{
		VenueId: venueID,
		OpeningHours: []*venue.OpeningHoursSpecification{
			{
				DayOfWeek: 1,
				Opens:     "10:00",
				Closes:    "22:00",
			},
		},
	}).Return(&api.UpdateOpeningHoursResponse{OpeningHours: []*venue.OpeningHoursSpecification{
		{
			DayOfWeek: 1,
			Opens:     "10:00",
			Closes:    "22:00",
		},
	}}, nil)

	var resp struct {
		UpdateOpeningHours []struct {
			DayOfWeek int    `json:"dayOfWeek"`
			Opens     string `json:"opens"`
			Closes    string `json:"closes"`
		} `json:"updateOpeningHours"`
	}

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	client.New(e).MustPost(`mutation{updateOpeningHours(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",openingHours:[{dayOfWeek:1,opens:"10:00",closes:"22:00"}]}) {dayOfWeek,opens,closes}}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_UpdateSpecialOpeningHours(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)
	date := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)
	venueClient.EXPECT().UpdateSpecialOpeningHours(gomock.Any(), &api.UpdateOpeningHoursRequest{
		VenueId: venueID,
		OpeningHours: []*venue.OpeningHoursSpecification{
			{
				DayOfWeek:    1,
				Opens:        "",
				Closes:       "",
				ValidFrom:    date,
				ValidThrough: date,
			},
		},
	}).Return(&api.UpdateOpeningHoursResponse{OpeningHours: []*venue.OpeningHoursSpecification{
		{
			DayOfWeek:    1,
			ValidFrom:    date,
			ValidThrough: date,
		},
	}}, nil)

	var resp struct {
		UpdateSpecialOpeningHours []struct {
			DayOfWeek    int    `json:"dayOfWeek"`
			Opens        string `json:"opens"`
			Closes       string `json:"closes"`
			ValidFrom    string `json:"validFrom"`
			ValidThrough string `json:"validThrough"`
		} `json:"updateSpecialOpeningHours"`
	}

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	client.New(e).MustPost(`mutation{updateSpecialOpeningHours(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",specialOpeningHours:[{dayOfWeek:1,validFrom:"3000-01-01T00:00:00Z",validThrough:"3000-01-01T00:00:00Z"}]}) {dayOfWeek,opens,closes,validFrom,validThrough}}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_RemoveTableNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		AddTable struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Capacity int    `json:"capacity"`
		} `json:"removeTable"`
	}
	assert.Error(t, c.Post(`mutation{removeTable(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",tableId:"bfcc0d78-83e7-4830-96ab-96cdbd0357c7"}) {id,name,capacity}}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_RemoveTable(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)
	venueClient.EXPECT().RemoveTable(gomock.Any(), &api.RemoveTableRequest{
		VenueId: venueID,
		TableId: "bfcc0d78-83e7-4830-96ab-96cdbd0357c7",
	}).Return(&venue.Table{
		Id:       "bfcc0d78-83e7-4830-96ab-96cdbd0357c7",
		Name:     "test table",
		Capacity: 5,
	}, nil)

	var resp struct {
		AddTable struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Capacity int    `json:"capacity"`
		} `json:"removeTable"`
	}

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	client.New(e).MustPost(`mutation{removeTable(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",tableId:"bfcc0d78-83e7-4830-96ab-96cdbd0357c7"}) {id,name,capacity}}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_AddAdminNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		AddAdmin string `json:"addAdmin"`
	}
	assert.Error(t, c.Post(`mutation{addAdmin(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",email:"test@test.com"})}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_AddAdmin(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)
	venueClient.EXPECT().AddAdmin(gomock.Any(), &api.AddAdminRequest{
		VenueId: venueID,
		Email:   "test@test.com",
	}).Return(&api.AddAdminResponse{VenueId: venueID, Email: "test@test.com"}, nil)

	var resp struct {
		AddAdmin string `json:"addAdmin"`
	}

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	client.New(e).MustPost(`mutation{addAdmin(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",email:"test@test.com"})}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_RemoveAdminNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		RemoveAdmin string `json:"removeAdmin"`
	}
	assert.Error(t, c.Post(`mutation{removeAdmin(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",email:"test@test.com"})}`, &resp), "user is not admin")
	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_RemoveAdmin(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)
	venueClient.EXPECT().RemoveAdmin(gomock.Any(), &api.RemoveAdminRequest{
		VenueId: venueID,
		Email:   "test@test.com",
	}).Return(&api.RemoveAdminResponse{Email: "test@test.com"}, nil)

	var resp struct {
		RemoveAdmin string `json:"removeAdmin"`
	}

	venueSrv, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueSrv, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	client.New(e).MustPost(`mutation{removeAdmin(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",email:"test@test.com"})}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_GetSlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingClient := mock_resolver.NewMockBookingAPIClient(ctrl)
	startsAt, err := time.Parse(time.RFC3339, "3000-06-20T12:41:45Z")
	require.NoError(t, err)

	bookingClient.EXPECT().GetSlot(gomock.Any(), &api2.SlotInput{
		VenueId:  "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:    "test@test.com",
		People:   5,
		StartsAt: startsAt.Format(time.RFC3339),
		Duration: 60,
	}).Return(&api2.GetSlotResponse{
		Match: &booking.Slot{
			VenueId:  "8a18e89b-339b-4e51-ab53-825aae59a070",
			Email:    "test@test.com",
			People:   5,
			StartsAt: startsAt.Format(time.RFC3339),
			EndsAt:   startsAt.Add(time.Minute * 60).Format(time.RFC3339),
			Duration: 60,
		},
		OtherAvailableSlots: nil,
	}, nil)

	bookingService, _, err := booking2.NewBookingClient("", nil, nil, booking2.WithClient(bookingClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), nil, bookingService)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		GetSlot struct {
			Match struct {
				VenueID  string `json:"venueId"`
				Email    string `json:"email"`
				People   int    `json:"people"`
				StartsAt string `json:"startsAt"`
				EndsAt   string `json:"endsAt"`
				Duration int    `json:"duration"`
			} `json:"match"`
		} `json:"getSlot"`
	}
	c.MustPost(`{getSlot(input:{venueId:"8a18e89b-339b-4e51-ab53-825aae59a070",email:"test@test.com",people:5,startsAt:"3000-06-20T12:41:45Z",duration:60,}) {match{venueId,email,people,startsAt,endsAt,duration}}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingClient := mock_resolver.NewMockBookingAPIClient(ctrl)
	startsAt, err := time.Parse(time.RFC3339, "3000-06-20T12:41:45Z")
	require.NoError(t, err)

	bookingClient.EXPECT().CreateBooking(gomock.Any(), &api2.BookingInput{
		VenueId:   "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:     "test@test.com",
		People:    5,
		StartsAt:  startsAt.Format(time.RFC3339),
		Duration:  60,
		Name:      "Test Test",
		GivenName: "Test",
	}).Return(&booking.Booking{
		Id:        "cca3c988-9e11-4b81-9a98-c960fb4a3d97",
		VenueId:   "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:     "test@test.com",
		People:    5,
		StartsAt:  startsAt.Format(time.RFC3339),
		EndsAt:    startsAt.Add(time.Minute * 60).Format(time.RFC3339),
		Duration:  60,
		TableId:   "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
		Name:      "Test Test",
		GivenName: "Test",
	}, nil)

	bookingService, _, err := booking2.NewBookingClient("", nil, nil, booking2.WithClient(bookingClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), nil, bookingService)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		CreateBooking struct {
			ID       string `json:"id"`
			VenueID  string `json:"venueId"`
			Email    string `json:"email"`
			People   int    `json:"people"`
			StartsAt string `json:"startsAt"`
			EndsAt   string `json:"endsAt"`
			Duration int    `json:"duration"`
			TableID  string `json:"tableId"`
		} `json:"createBooking"`
	}
	c.MustPost(`mutation{createBooking(input:{venueId:"8a18e89b-339b-4e51-ab53-825aae59a070",email:"test@test.com",people:5,startsAt:"3000-06-20T12:41:45Z",duration:60,}) {id,venueId,email,people,startsAt,endsAt,duration,tableId}}`, &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_IsAdminTrue(t *testing.T) {
	var venueID = "8a18e89b-339b-4e51-ab53-825aae59a070"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)

	venueService, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueService, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		IsAdmin bool `json:"isAdmin"`
	}
	c.MustPost(fmt.Sprintf(`{isAdmin(input:{venueId:"%s"})}`, venueID), &resp)

	if resp.IsAdmin != true {
		t.Errorf("expected is admin == true, got false")
	}

	ctrl.Finish()
}

func Test_IsAdminFalse(t *testing.T) {
	var venueID = "8a18e89b-339b-4e51-ab53-825aae59a070"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueService, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueService, nil)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		IsAdmin bool `json:"isAdmin"`
	}
	c.MustPost(fmt.Sprintf(`{isAdmin(input:{venueId:"%s"})}`, venueID), &resp)

	if resp.IsAdmin != false {
		t.Errorf("expected is admin == false, got true")
	}

	ctrl.Finish()
}

func Test_CancelBooking(t *testing.T) {
	venueID := "8a18e89b-339b-4e51-ab53-825aae59a070"
	bookingID := "47f4eaf4-7b5e-43dc-bc06-ebf8561c1fa9"
	startsAt, err := time.Parse(time.RFC3339, "3000-06-20T12:41:45Z")
	require.NoError(t, err)
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)
	bookingClient := mock_resolver.NewMockBookingAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: true}, nil)

	bookingClient.EXPECT().CancelBooking(gomock.Any(), &api2.CancelBookingRequest{
		Id: bookingID,
	}).Return(&booking.Booking{
		Id:        "cca3c988-9e11-4b81-9a98-c960fb4a3d97",
		VenueId:   "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:     "test@test.com",
		People:    5,
		StartsAt:  startsAt.Format(time.RFC3339),
		EndsAt:    startsAt.Add(time.Minute * 60).Format(time.RFC3339),
		Duration:  60,
		TableId:   "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
		Name:      "Test Test",
		GivenName: "Test",
	}, nil)

	venueService, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)
	bookingService, _, err := booking2.NewBookingClient("", nil, nil, booking2.WithClient(bookingClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueService, bookingService)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		CancelBooking struct {
			ID       string `json:"id"`
			VenueID  string `json:"venueId"`
			Email    string `json:"email"`
			People   int    `json:"people"`
			StartsAt string `json:"startsAt"`
			EndsAt   string `json:"endsAt"`
			Duration int    `json:"duration"`
			TableID  string `json:"tableId"`
		} `json:"cancelBooking"`
	}
	c.MustPost(fmt.Sprintf(`mutation{cancelBooking(input:{venueId:"%s",id:"%s"}){id,venueId,email,people,startsAt,endsAt,duration,tableId}}`, venueID, bookingID), &resp)

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func Test_CancelBookingNotAuthorised(t *testing.T) {
	venueID := "8a18e89b-339b-4e51-ab53-825aae59a070"
	bookingID := "47f4eaf4-7b5e-43dc-bc06-ebf8561c1fa9"
	ctrl := gomock.NewController(t)
	venueClient := mock_resolver.NewMockVenueAPIClient(ctrl)
	bookingClient := mock_resolver.NewMockBookingAPIClient(ctrl)

	venueClient.EXPECT().IsAdmin(gomock.Any(), &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    "",
		Email:   "test@test.com",
	}).Return(&api.IsAdminResponse{IsAdmin: false}, nil)

	venueService, _, err := venue2.NewVenueClient("", nil, nil, venue2.WithClient(venueClient))
	require.NoError(t, err)
	bookingService, _, err := booking2.NewBookingClient("", nil, nil, booking2.WithClient(bookingClient))
	require.NoError(t, err)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(zap.NewNop().Sugar(), venueService, bookingService)}))
	e := echo.New()
	e.POST("/", echo.WrapHandler(h), middleware.User(mockUserService{}))
	c := client.New(e)

	var resp struct {
		CancelBooking struct {
			ID       string `json:"id"`
			VenueID  string `json:"venueId"`
			Email    string `json:"email"`
			People   int    `json:"people"`
			StartsAt string `json:"startsAt"`
			EndsAt   string `json:"endsAt"`
			Duration int    `json:"duration"`
			TableID  string `json:"tableId"`
		} `json:"cancelBooking"`
	}
	assert.Error(t, c.Post(fmt.Sprintf(`mutation{cancelBooking(input:{venueId:"%s",id:"%s"}){id,venueId,email,people,startsAt,endsAt,duration,tableId}}`, venueID, bookingID), &resp), "user is not admin")

	cupaloy.SnapshotT(t, resp)

	ctrl.Finish()
}

func defaultOpeningHours() []*venue.OpeningHoursSpecification {
	return []*venue.OpeningHoursSpecification{{
		DayOfWeek:    1,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    "",
		ValidThrough: "",
	}, {
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    "",
		ValidThrough: "",
	}}
}

var _ models.UserService = (*mockUserService)(nil)

type mockUserService struct{}

func (m mockUserService) GetUser(ctx context.Context) (*models.User, error) {
	return &models.User{
		Name:      "Test Test",
		Email:     "test@test.com",
		GivenName: "Test",
	}, nil
}
