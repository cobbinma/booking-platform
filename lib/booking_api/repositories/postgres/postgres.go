package postgres

import (
	"github.com/cobbinma/booking/lib/booking_api/models"
)

type postgres struct {
	dbClient DBClient
}

func NewPostgres(client DBClient) models.Repository {
	return &postgres{dbClient: client}
}
