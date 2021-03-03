package graph_test

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	mock_resolver "github.com/cobbinma/booking-platform/lib/gateway_api/graph/mock"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_GetVenue(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueSrv := mock_resolver.NewMockVenueService(ctrl)
	monday := &models.OpeningHoursSpecification{
		DayOfWeek:    models.Monday,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}
	tuesday := &models.OpeningHoursSpecification{
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}

	venueSrv.EXPECT().GetVenue(gomock.Any(), models.VenueFilter{
		ID: &venueID,
	}).Return(&models.Venue{
		ID:                  "a3291740-e89f-4cc0-845c-75c4c39842c9",
		Name:                "hop and vine",
		OpeningHours:        []*models.OpeningHoursSpecification{monday, tuesday},
		SpecialOpeningHours: nil,
		Slug:                "hop-and-vine",
	}, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(nil, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)
	monday := &models.OpeningHoursSpecification{
		DayOfWeek:    models.Monday,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}
	tuesday := &models.OpeningHoursSpecification{
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}

	venueSrv.EXPECT().GetVenue(gomock.Any(), models.VenueFilter{
		Slug: &slug,
	}).Return(&models.Venue{
		ID:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        []*models.OpeningHoursSpecification{monday, tuesday},
		SpecialOpeningHours: nil,
	}, nil)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)

	venueSrv.EXPECT().GetTables(gomock.Any(), venueID).Return([]*models.Table{
		{
			ID:       "175fd06d-9a60-4ea6-86ca-bb96ca861208",
			Name:     "table one",
			Capacity: 4,
		},
	}, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)
	monday := &models.OpeningHoursSpecification{
		DayOfWeek:    models.Monday,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}
	tuesday := &models.OpeningHoursSpecification{
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}

	venueSrv.EXPECT().GetVenue(gomock.Any(), models.VenueFilter{
		ID: &venueID,
	}).Return(&models.Venue{
		ID:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        []*models.OpeningHoursSpecification{monday, tuesday},
		SpecialOpeningHours: nil,
	}, nil)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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

func Test_GetVenueAdmins(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	slug := "test-venue"
	ctrl := gomock.NewController(t)
	venueSrv := mock_resolver.NewMockVenueService(ctrl)
	monday := &models.OpeningHoursSpecification{
		DayOfWeek:    models.Monday,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}
	tuesday := &models.OpeningHoursSpecification{
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}

	venueSrv.EXPECT().GetVenue(gomock.Any(), models.VenueFilter{
		Slug: &slug,
	}).Return(&models.Venue{
		ID:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        []*models.OpeningHoursSpecification{monday, tuesday},
		SpecialOpeningHours: nil,
	}, nil)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)

	venueSrv.EXPECT().GetAdmins(gomock.Any(), venueID).Return([]string{"test@test.com"}, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)
	monday := &models.OpeningHoursSpecification{
		DayOfWeek:    models.Monday,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}
	tuesday := &models.OpeningHoursSpecification{
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    nil,
		ValidThrough: nil,
	}

	venueSrv.EXPECT().GetVenue(gomock.Any(), models.VenueFilter{
		ID: &venueID,
	}).Return(&models.Venue{
		ID:                  venueID,
		Name:                "hop and vine",
		OpeningHours:        []*models.OpeningHoursSpecification{monday, tuesday},
		SpecialOpeningHours: nil,
	}, nil)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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

func Test_AddTableNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)
	venueSrv.EXPECT().AddTable(gomock.Any(), models.TableInput{
		VenueID:  venueID,
		Name:     "test table",
		Capacity: 5,
	}).Return(&models.Table{
		ID:       "bfcc0d78-83e7-4830-96ab-96cdbd0357c7",
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
	client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)}))).
		MustPost(`mutation{addTable(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",name:"test table",capacity:5}) {id,name,capacity}}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_RemoveTableNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)
	venueSrv.EXPECT().RemoveTable(gomock.Any(), models.RemoveTableInput{
		VenueID: venueID,
		TableID: "bfcc0d78-83e7-4830-96ab-96cdbd0357c7",
	}).Return(&models.Table{
		ID:       "bfcc0d78-83e7-4830-96ab-96cdbd0357c7",
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
	client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)}))).
		MustPost(`mutation{removeTable(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",tableId:"bfcc0d78-83e7-4830-96ab-96cdbd0357c7"}) {id,name,capacity}}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_AddAdminNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)
	venueSrv.EXPECT().AddAdmin(gomock.Any(), models.AdminInput{
		VenueID: venueID,
		Email:   "test@test.com",
	}).Return("test@test.com", nil)

	var resp struct {
		AddAdmin string `json:"addAdmin"`
	}
	client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)}))).
		MustPost(`mutation{addAdmin(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",email:"test@test.com"})}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_RemoveAdminNotAuthorised(t *testing.T) {
	venueID := "a3291740-e89f-4cc0-845c-75c4c39842c9"
	ctrl := gomock.NewController(t)
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)})))

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
	venueSrv := mock_resolver.NewMockVenueService(ctrl)

	venueSrv.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)
	venueSrv.EXPECT().RemoveAdmin(gomock.Any(), models.RemoveAdminInput{
		VenueID: venueID,
		Email:   "test@test.com",
	}).Return("test@test.com", nil)

	var resp struct {
		RemoveAdmin string `json:"removeAdmin"`
	}
	client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueSrv, nil)}))).
		MustPost(`mutation{removeAdmin(input:{venueId:"a3291740-e89f-4cc0-845c-75c4c39842c9",email:"test@test.com"})}`, &resp)

	cupaloy.SnapshotT(t, resp)
	ctrl.Finish()
}

func Test_GetSlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookingService := mock_resolver.NewMockBookingService(ctrl)
	startsAt, err := time.Parse(time.RFC3339, "3000-06-20T12:41:45Z")
	require.NoError(t, err)

	bookingService.EXPECT().GetSlot(gomock.Any(), models.SlotInput{
		VenueID:  "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:    "test@test.com",
		People:   5,
		StartsAt: startsAt,
		Duration: 60,
	}).Return(&models.GetSlotResponse{
		Match: &models.Slot{
			VenueID:  "8a18e89b-339b-4e51-ab53-825aae59a070",
			Email:    "test@test.com",
			People:   5,
			StartsAt: startsAt,
			EndsAt:   startsAt.Add(time.Minute * 60),
			Duration: 60,
		},
		OtherAvailableSlots: nil,
	}, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, nil, bookingService)})))

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
	bookingService := mock_resolver.NewMockBookingService(ctrl)
	startsAt, err := time.Parse(time.RFC3339, "3000-06-20T12:41:45Z")
	require.NoError(t, err)

	bookingService.EXPECT().CreateBooking(gomock.Any(), models.BookingInput{
		VenueID:  "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:    "test@test.com",
		People:   5,
		StartsAt: startsAt,
		Duration: 60,
	}).Return(&models.Booking{
		ID:       "cca3c988-9e11-4b81-9a98-c960fb4a3d97",
		VenueID:  "8a18e89b-339b-4e51-ab53-825aae59a070",
		Email:    "test@test.com",
		People:   5,
		StartsAt: startsAt,
		EndsAt:   startsAt.Add(time.Minute * 60),
		Duration: 60,
		TableID:  "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
	}, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, nil, bookingService)})))

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
	venueService := mock_resolver.NewMockVenueService(ctrl)

	venueService.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(true, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueService, nil)})))

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
	venueService := mock_resolver.NewMockVenueService(ctrl)

	venueService.EXPECT().IsAdmin(gomock.Any(), models.IsAdminInput{VenueID: &venueID}, "test@test.com").Return(false, nil)

	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, venueService, nil)})))

	var resp struct {
		IsAdmin bool `json:"isAdmin"`
	}
	c.MustPost(fmt.Sprintf(`{isAdmin(input:{venueId:"%s"})}`, venueID), &resp)

	if resp.IsAdmin != false {
		t.Errorf("expected is admin == false, got true")
	}

	ctrl.Finish()
}

var _ models.UserService = (*mockUserService)(nil)

type mockUserService struct{}

func (m mockUserService) GetUser(ctx context.Context) (*models.User, error) {
	return &models.User{
		Name:  "Test Test",
		Email: "test@test.com",
	}, nil
}
