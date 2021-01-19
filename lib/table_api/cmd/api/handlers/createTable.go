package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func CreateTable(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		reqBody, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not read request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
		}

		table := models.NewTable{}
		if err := json.Unmarshal(reqBody, &table); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not unmarshall request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
		}

		if err := table.Valid(); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "invalid request", err))
			message := fmt.Sprintf("incorrect user request : %s", err)
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
		}

		tbl, err := repository.CreateTable(ctx, table)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not create table", err))
			message := "could not create table"
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
		}

		return c.JSON(http.StatusCreated, tbl)
	}
}
