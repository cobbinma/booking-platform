package postgres

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var _ api.VenueAPIServer = (*client)(nil)
var _ api.TableAPIServer = (*client)(nil)

type client struct {
	db  *sqlx.DB
	log *zap.SugaredLogger
}

func NewPostgres(log *zap.SugaredLogger, options ...func(*client)) (*client, func(log *zap.SugaredLogger), error) {
	c := &client{log: log}
	for i := range options {
		options[i](c)
	}

	return c, nil, nil
}

func (c client) GetTables(ctx context.Context, request *api.GetTablesRequest) (*api.GetTablesResponse, error) {
	panic("implement me")
}

func (c client) AddTable(ctx context.Context, request *api.AddTableRequest) (*models.Table, error) {
	panic("implement me")
}

func (c client) GetVenue(ctx context.Context, request *api.GetVenueRequest) (*models.Venue, error) {
	panic("implement me")
}

func (c client) CreateVenue(ctx context.Context, request *api.CreateVenueRequest) (*models.Venue, error) {
	panic("implement me")
}
