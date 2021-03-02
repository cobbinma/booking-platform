package venue

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"time"
)

func NewVenueClient(url string, log *zap.SugaredLogger, token *oauth2.Token) (graph.VenueService, func(log *zap.SugaredLogger), error) {
	c, err := credentials.NewClientTLSFromFile("localhost.crt", "localhost")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load credentials : %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)),
		grpc.WithTransportCredentials(c),
	}
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect : %s", err)
	}

	return &venueClient{
			client: api.NewVenueAPIClient(conn),
			log:    log,
		}, func(log *zap.SugaredLogger) {
			if err := conn.Close(); err != nil {
				log.Error("could not close connection : %s", err)
			}
		}, nil
}

type venueClient struct {
	client api.VenueAPIClient
	log    *zap.SugaredLogger
}

func (v venueClient) AddTable(ctx context.Context, input models.TableInput) (*models.Table, error) {
	table, err := v.client.AddTable(ctx, &api.AddTableRequest{
		VenueId:  input.VenueID,
		Name:     input.Name,
		Capacity: uint32(input.Capacity),
	})
	if err != nil {
		return nil, fmt.Errorf("could not add table using venue service : %w", err)
	}

	return &models.Table{
		ID:       table.Id,
		Name:     table.Name,
		Capacity: int(table.Capacity),
	}, nil
}

func (v venueClient) RemoveTable(ctx context.Context, input models.RemoveTableInput) (*models.Table, error) {
	table, err := v.client.RemoveTable(ctx, &api.RemoveTableRequest{
		VenueId: input.VenueID,
		TableId: input.TableID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not add table using venue service : %w", err)
	}

	return &models.Table{
		ID:       table.Id,
		Name:     table.Name,
		Capacity: int(table.Capacity),
	}, nil
}

func (v venueClient) GetTables(ctx context.Context, venueID string) ([]*models.Table, error) {
	resp, err := v.client.GetTables(ctx, &api.GetTablesRequest{VenueId: venueID})
	if err != nil {
		return nil, fmt.Errorf("could not get tables from venue service : %w", err)
	}

	tables := []*models.Table{}
	for _, table := range resp.Tables {
		tables = append(tables, &models.Table{
			ID:       table.Id,
			Name:     table.Name,
			Capacity: int(table.Capacity),
		})
	}

	return tables, nil
}

func (v venueClient) GetVenue(ctx context.Context, filter models.VenueFilter) (*models.Venue, error) {
	var id, slug string
	if filter.ID != nil {
		id = *filter.ID
	}
	if filter.Slug != nil {
		slug = *filter.Slug
	}
	venue, err := v.client.GetVenue(ctx, &api.GetVenueRequest{Id: id, Slug: slug})
	if err != nil {
		return nil, fmt.Errorf("could not get venue from venue service : %w", err)
	}

	openingHours := []*models.OpeningHoursSpecification{}
	for _, hours := range venue.OpeningHours {
		openingHours = append(openingHours, &models.OpeningHoursSpecification{
			DayOfWeek:    (models.DayOfWeek)(hours.DayOfWeek),
			Opens:        (models.TimeOfDay)(hours.Opens),
			Closes:       (models.TimeOfDay)(hours.Closes),
			ValidFrom:    nil,
			ValidThrough: nil,
		})
	}

	specialHours := []*models.OpeningHoursSpecification{}
	for _, hours := range venue.SpecialOpeningHours {
		from, err := time.Parse(time.RFC3339, hours.ValidFrom)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid from '%s'", err)
		}
		through, err := time.Parse(time.RFC3339, hours.ValidThrough)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid through '%s'", err)
		}
		openingHours = append(specialHours, &models.OpeningHoursSpecification{
			DayOfWeek:    (models.DayOfWeek)(hours.DayOfWeek),
			Opens:        (models.TimeOfDay)(hours.Opens),
			Closes:       (models.TimeOfDay)(hours.Closes),
			ValidFrom:    &from,
			ValidThrough: &through,
		})
	}

	return &models.Venue{
		ID:                  venue.Id,
		Name:                venue.Name,
		OpeningHours:        openingHours,
		SpecialOpeningHours: specialHours,
		Slug:                venue.Slug,
	}, nil
}

func (v venueClient) IsAdmin(ctx context.Context, venueID string, email string) (bool, error) {
	resp, err := v.client.IsAdmin(ctx, &api.IsAdminRequest{
		VenueId: venueID,
		Email:   email,
	})
	if err != nil {
		return false, fmt.Errorf("could not get is admin from customer client : %w", err)
	}

	return resp.IsAdmin, nil
}
