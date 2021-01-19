package handlers_test

import (
	"context"
	"encoding/json"
	"github.com/cobbinma/booking-platform/lib/venue_api/cmd/api/handlers"
	"github.com/cobbinma/booking-platform/lib/venue_api/models"
	"github.com/cobbinma/booking-platform/lib/venue_api/repositories/fakeRepository"
	"github.com/cobbinma/booking-platform/lib/venue_api/repositories/postgres"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeleteVenueNotFound(t *testing.T) {
	repository := fakeRepository.NewFakeRepository()

	id := 1
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	deleteVenue := handlers.DeleteVenue(repository)

	if err := deleteVenue(c); err != nil {
		t.Errorf("error returned from delete venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusNoContent)
		return
	}
}

func TestDeleteVenueFound(t *testing.T) {
	repository := fakeRepository.NewFakeRepository()

	var input models.VenueInput
	if err := json.Unmarshal([]byte(`{"name":"hop and vine","openingHours":[{"dayOfWeek":1,"opens":"09:00","closes":"22:00"}]}`), &input); err != nil {
		t.Errorf("could not unmarshall venue json : %s", err)
		return
	}

	if _, err := repository.CreateVenue(context.Background(), input); err != nil {
		t.Errorf("could not create venue in repository : %s", err)
		return
	}

	id := 0
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	deleteVenue := handlers.DeleteVenue(repository)

	if err := deleteVenue(c); err != nil {
		t.Errorf("error returned from delete venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusNoContent)
		return
	}

	_, err := repository.GetVenue(context.Background(), models.VenueID(id))
	if !postgres.ErrVenueNotFound(err) {
		t.Errorf("should have returned error venue not found")
		return
	}
}
