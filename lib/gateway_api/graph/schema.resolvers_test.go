package graph_test

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	venue "github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"google.golang.org/grpc"
	"testing"
)

func Test_GetVenue(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, &mockVenueService{})})))

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
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(id:"a3291740-e89f-4cc0-845c-75c4c39842c9"){id,name,openingHours{dayOfWeek,opens,closes,validFrom,validThrough},specialOpeningHours{dayOfWeek,opens, closes, validFrom,validThrough}}}`, &resp)

	cupaloy.SnapshotT(t, resp)
}

func Test_GetSlot(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, &mockVenueService{})})))

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
}

func Test_CreateBooking(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{}, &mockVenueService{})})))

	var resp struct {
		CreateBooking struct {
			VenueID  string `json:"venueId"`
			Email    string `json:"email"`
			People   int    `json:"people"`
			StartsAt string `json:"startsAt"`
			EndsAt   string `json:"endsAt"`
			Duration int    `json:"duration"`
			TableID  string `json:"tableId"`
		} `json:"createBooking"`
	}
	c.MustPost(`mutation{createBooking(input:{venueId:"8a18e89b-339b-4e51-ab53-825aae59a070",email:"test@test.com",people:5,startsAt:"3000-06-20T12:41:45Z",duration:60,}) {venueId,email,people,startsAt,endsAt,duration,tableId}}`, &resp)

	cupaloy.SnapshotT(t, resp)
}

var _ models.UserService = (*mockUserService)(nil)

type mockUserService struct{}

func (m mockUserService) GetUser(ctx context.Context) (*models.User, error) {
	return &models.User{
		Name:  "Test Test",
		Email: "test@test.com",
	}, nil
}

var _ api.VenueAPIClient = (*mockVenueService)(nil)

type mockVenueService struct{}

func (m mockVenueService) GetVenue(ctx context.Context, in *api.GetVenueRequest, opts ...grpc.CallOption) (*venue.Venue, error) {
	monday := &venue.OpeningHoursSpecification{
		DayOfWeek:    1,
		Opens:        "10:00",
		Closes:       "19:00",
		ValidFrom:    "",
		ValidThrough: "",
	}
	tuesday := &venue.OpeningHoursSpecification{
		DayOfWeek:    2,
		Opens:        "11:00",
		Closes:       "20:00",
		ValidFrom:    "",
		ValidThrough: "",
	}
	return &venue.Venue{
		Id:                  in.Id,
		Name:                "hop and vine",
		OpeningHours:        []*venue.OpeningHoursSpecification{monday, tuesday},
		SpecialOpeningHours: nil,
	}, nil
}
