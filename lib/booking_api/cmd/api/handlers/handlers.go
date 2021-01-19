package handlers

import "github.com/cobbinma/booking-platform/lib/booking_api/models"

type Handlers struct {
	repository  models.Repository
	tableClient models.TableClient
}

func NewHandlers(repository models.Repository, tableClient models.TableClient) *Handlers {
	return &Handlers{repository: repository, tableClient: tableClient}
}
