package postgres

import "github.com/cobbinma/booking/lib/venue_api/models"

var _ models.Repository = (*postgres)(nil)

type postgres struct {
	dbClient DBClient
}

func NewPostgres(client DBClient) models.Repository {
	return &postgres{dbClient: client}
}
