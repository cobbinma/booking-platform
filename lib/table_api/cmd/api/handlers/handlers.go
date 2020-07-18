package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Handlers struct {
	repository models.Repository
}

func NewHandlers(repository models.Repository) *Handlers {
	return &Handlers{repository: repository}
}

func getTableIDFromRequest(c echo.Context) (models.TableID, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not parse id", err)
	}

	tableID := models.NewTableID(id)
	err = tableID.Valid()
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not validate id", err)
	}

	return tableID, nil
}
