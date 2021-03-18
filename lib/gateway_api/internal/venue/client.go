package venue

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	venue "github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
	"time"
)

func NewVenueClient(url string, log *zap.SugaredLogger, token *oauth2.Token, options ...func(*venueClient)) (graph.VenueService, func(log *zap.SugaredLogger), error) {
	vc := &venueClient{
		client: nil,
		log:    log,
	}
	cl := func(log *zap.SugaredLogger) {}

	for i := range options {
		options[i](vc)
	}

	if vc.client == nil {
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

		cl = func(log *zap.SugaredLogger) {
			if err := conn.Close(); err != nil {
				log.Error("could not close connection : %s", err)
			}
		}
		vc.client = api.NewVenueAPIClient(conn)
	}

	return vc, cl, nil
}

func WithClient(client api.VenueAPIClient) func(*venueClient) {
	return func(c *venueClient) {
		c.client = client
	}
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
			Opens:        (*models.TimeOfDay)(&hours.Opens),
			Closes:       (*models.TimeOfDay)(&hours.Closes),
			ValidFrom:    nil,
			ValidThrough: nil,
		})
	}

	specialHours := []*models.OpeningHoursSpecification{}
	for _, hours := range venue.SpecialOpeningHours {
		var opens, closes *models.TimeOfDay
		if hours.Opens != "" {
			opens = (*models.TimeOfDay)(&hours.Opens)
		}
		if hours.Closes != "" {
			closes = (*models.TimeOfDay)(&hours.Closes)
		}
		from, err := time.Parse(time.RFC3339, hours.ValidFrom)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid from '%s'", err)
		}
		through, err := time.Parse(time.RFC3339, hours.ValidThrough)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid through '%s'", err)
		}
		specialHours = append(specialHours, &models.OpeningHoursSpecification{
			DayOfWeek:    (models.DayOfWeek)(hours.DayOfWeek),
			Opens:        opens,
			Closes:       closes,
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

func (v venueClient) OpeningHoursSpecification(ctx context.Context, venueID string, date time.Time) (*models.OpeningHoursSpecification, error) {
	resp, err := v.client.GetOpeningHoursSpecification(ctx, &api.GetOpeningHoursSpecificationRequest{
		VenueId: venueID,
		Date:    date.Format(time.RFC3339),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not get specification from client : %w", err)
	}

	var opens, closes *models.TimeOfDay
	if resp.Specification.Opens != "" {
		opens = (*models.TimeOfDay)(&resp.Specification.Opens)
	}
	if resp.Specification.Closes != "" {
		closes = (*models.TimeOfDay)(&resp.Specification.Closes)
	}

	var validFrom, validThrough *time.Time
	if resp.Specification.ValidFrom != "" && resp.Specification.ValidThrough != "" {
		f, err := time.Parse(time.RFC3339, resp.Specification.ValidFrom)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid from time : %w", err)
		}
		validFrom = &f

		t, err := time.Parse(time.RFC3339, resp.Specification.ValidThrough)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid through time : %w", err)
		}
		validThrough = &t
	}

	return &models.OpeningHoursSpecification{
		DayOfWeek:    models.DayOfWeek(resp.Specification.DayOfWeek),
		Opens:        opens,
		Closes:       closes,
		ValidFrom:    validFrom,
		ValidThrough: validThrough,
	}, nil
}

func (v venueClient) UpdateOpeningHours(ctx context.Context, input models.UpdateOpeningHoursInput) ([]*models.OpeningHoursSpecification, error) {
	hours := make([]*venue.OpeningHoursSpecification, len(input.OpeningHours))
	for i := range input.OpeningHours {
		hours[i] = &venue.OpeningHoursSpecification{
			DayOfWeek: uint32(input.OpeningHours[i].DayOfWeek),
			Opens:     string(input.OpeningHours[i].Opens),
			Closes:    string(input.OpeningHours[i].Closes),
		}
	}

	resp, err := v.client.UpdateOpeningHours(ctx, &api.UpdateOpeningHoursRequest{
		VenueId:      input.VenueID,
		OpeningHours: hours,
	})
	if err != nil {
		return nil, fmt.Errorf("could not update opening hours in client : %w", err)
	}

	updated := make([]*models.OpeningHoursSpecification, len(resp.OpeningHours))
	for i := range resp.OpeningHours {
		opens := models.TimeOfDay(resp.OpeningHours[i].Opens)
		closes := models.TimeOfDay(resp.OpeningHours[i].Closes)
		updated[i] = &models.OpeningHoursSpecification{
			DayOfWeek: models.DayOfWeek(resp.OpeningHours[i].DayOfWeek),
			Opens:     &opens,
			Closes:    &closes,
		}
	}

	return updated, nil
}

func (v venueClient) UpdateSpecialOpeningHours(ctx context.Context, input models.UpdateSpecialOpeningHoursInput) ([]*models.OpeningHoursSpecification, error) {
	hours := make([]*venue.OpeningHoursSpecification, len(input.SpecialOpeningHours))
	for i := range input.SpecialOpeningHours {
		var opens, closes string
		if input.SpecialOpeningHours[i].Opens != nil {
			opens = string(*input.SpecialOpeningHours[i].Opens)
		}
		if input.SpecialOpeningHours[i].Closes != nil {
			closes = string(*input.SpecialOpeningHours[i].Closes)
		}
		hours[i] = &venue.OpeningHoursSpecification{
			DayOfWeek:    uint32(input.SpecialOpeningHours[i].DayOfWeek),
			Opens:        opens,
			Closes:       closes,
			ValidFrom:    input.SpecialOpeningHours[i].ValidFrom.Format(time.RFC3339),
			ValidThrough: input.SpecialOpeningHours[i].ValidThrough.Format(time.RFC3339),
		}
	}

	resp, err := v.client.UpdateSpecialOpeningHours(ctx, &api.UpdateOpeningHoursRequest{
		VenueId:      input.VenueID,
		OpeningHours: hours,
	})
	if err != nil {
		return nil, fmt.Errorf("could not update special opening hours in client : %w", err)
	}

	updated := make([]*models.OpeningHoursSpecification, len(resp.OpeningHours))
	for i := range resp.OpeningHours {
		var opens, closes *models.TimeOfDay
		if resp.OpeningHours[i].Opens != "" {
			opens = (*models.TimeOfDay)(&resp.OpeningHours[i].Opens)
		}
		if resp.OpeningHours[i].Closes != "" {
			closes = (*models.TimeOfDay)(&resp.OpeningHours[i].Closes)
		}
		from, err := time.Parse(time.RFC3339, resp.OpeningHours[i].ValidFrom)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid from : %w", err)
		}
		through, err := time.Parse(time.RFC3339, resp.OpeningHours[i].ValidThrough)
		if err != nil {
			return nil, fmt.Errorf("could not parse valid through : %w", err)
		}
		updated[i] = &models.OpeningHoursSpecification{
			DayOfWeek:    models.DayOfWeek(resp.OpeningHours[i].DayOfWeek),
			Opens:        opens,
			Closes:       closes,
			ValidFrom:    &from,
			ValidThrough: &through,
		}
	}

	return updated, nil
}

func (v venueClient) IsAdmin(ctx context.Context, input models.IsAdminInput, email string) (bool, error) {
	var venueID, slug string
	if input.VenueID != nil {
		venueID = *input.VenueID
	}
	if input.Slug != nil {
		slug = *input.Slug
	}
	resp, err := v.client.IsAdmin(ctx, &api.IsAdminRequest{
		VenueId: venueID,
		Slug:    slug,
		Email:   email,
	})
	if err != nil {
		return false, fmt.Errorf("could not get is admin from client : %w", err)
	}

	return resp.IsAdmin, nil
}

func (v venueClient) GetAdmins(ctx context.Context, venueID string) ([]string, error) {
	resp, err := v.client.GetAdmins(ctx, &api.GetAdminsRequest{VenueId: venueID})
	if err != nil {
		return nil, fmt.Errorf("could not get admins from client : %w", err)
	}

	return resp.Admins, nil
}

func (v venueClient) AddAdmin(ctx context.Context, input models.AdminInput) (string, error) {
	resp, err := v.client.AddAdmin(ctx, &api.AddAdminRequest{
		VenueId: input.VenueID,
		Email:   input.Email,
	})
	if err != nil {
		return "", fmt.Errorf("could not add admin using client : %w", err)
	}

	return resp.Email, nil
}

func (v venueClient) RemoveAdmin(ctx context.Context, input models.RemoveAdminInput) (string, error) {
	resp, err := v.client.RemoveAdmin(ctx, &api.RemoveAdminRequest{
		VenueId: input.VenueID,
		Email:   input.Email,
	})
	if err != nil {
		return "", fmt.Errorf("could not remove admin using client : %w", err)
	}

	return resp.Email, nil
}
