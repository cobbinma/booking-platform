package graph_test

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
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
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"openingHours"`
			SpecialOpeningHours []struct {
				DayOfWeek    int    `json:"dayOfWeek"`
				ValidFrom    string `json:"validFrom"`
				ValidThrough string `json:"validThrough"`
			} `json:"specialOpeningHours"`
		} `json:"getVenue"`
	}
	c.MustPost(`{getVenue(id:"a3291740-e89f-4cc0-845c-75c4c39842c9"){id,name,openingHours{dayOfWeek,validFrom,validThrough},specialOpeningHours{dayOfWeek,validFrom,validThrough}}}`, &resp)

	cupaloy.SnapshotT(t, resp)
}
