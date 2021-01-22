package graph_test

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"testing"
)

func Test_GetVenue(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})))

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
	c.MustPost(`{getVenue(id:"a3291740-e89f-4cc0-845c-75c4c39842c9"){id,name,openingHours{dayOfWeek,validFrom,validThrough},specialOpeningHours{dayOfWeek,validFrom,validThrough}}}`, &resp)

	cupaloy.SnapshotT(t, resp)
}

func Test_CreateSlot(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})))

	var resp struct {
		CreateSlot struct {
			VenueID    string `json:"venueId"`
			CustomerID string `json:"customerId"`
			People     int    `json:"people"`
			Date       string `json:"date"`
			StartsAt   string `json:"startsAt"`
			EndsAt     string `json:"endsAt"`
			Duration   int    `json:"duration"`
		} `json:"createSlot"`
	}
	c.MustPost(`mutation{createSlot(input:{venueId:"8a18e89b-339b-4e51-ab53-825aae59a070",customerId:"23a4d31c-d6e4-4cfc-9cf0-b00b08faba55",people:5,date:"01-05-3000",startsAt:"15:00",duration:60,}) {venueId,customerId,people,date,startsAt,endsAt,duration}}`, &resp)

	cupaloy.SnapshotT(t, resp)
}

func Test_CreateBooking(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(&mockUserService{})})))

	var resp struct {
		CreateBooking struct {
			VenueID    string `json:"venueId"`
			CustomerID string `json:"customerId"`
			People     int    `json:"people"`
			Date       string `json:"date"`
			StartsAt   string `json:"startsAt"`
			EndsAt     string `json:"endsAt"`
			Duration   int    `json:"duration"`
			TableID    string `json:"tableId"`
		} `json:"createBooking"`
	}
	c.MustPost(`mutation{createBooking(input:{venueId:"8a18e89b-339b-4e51-ab53-825aae59a070",customerId:"23a4d31c-d6e4-4cfc-9cf0-b00b08faba55",people:5,date:"01-05-3000",startsAt:"15:00",duration:60,}) {venueId,customerId,people,date,startsAt,endsAt,duration,tableId}}`, &resp)

	cupaloy.SnapshotT(t, resp)
}

var _ models.UserService = (*mockUserService)(nil)

type mockUserService struct{}

func (m mockUserService) GetUser(ctx context.Context) (*models.User, error) {
	return &models.User{
		Name:  "Test Test",
		Email: "23a4d31c-d6e4-4cfc-9cf0-b00b08faba55",
	}, nil
}
