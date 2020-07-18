package handlers

import "github.com/cobbinma/booking/lib/booking_api/models"

type Handlers struct {
	repository models.Repository
}

func NewHandlers(repository models.Repository) *Handlers {
	return &Handlers{repository: repository}
}
