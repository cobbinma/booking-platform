package handlers_test

import (
	"context"
	"encoding/json"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/venue_api/cmd/api/handlers"
	"github.com/cobbinma/booking-platform/lib/venue_api/models"
	"github.com/cobbinma/booking-platform/lib/venue_api/repositories/fakeRepository"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetVenueNotFound(t *testing.T) {
	repository := fakeRepository.NewFakeRepository()

	id := 1
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	getVenue := handlers.GetVenue(repository)

	if err := getVenue(c); err != nil {
		t.Errorf("error returned from get venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusNotFound)
		return
	}

	cupaloy.SnapshotT(t, rec.Body.String())
}

func TestGetVenueFound(t *testing.T) {
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

	getVenue := handlers.GetVenue(repository)

	if err := getVenue(c); err != nil {
		t.Errorf("error returned from get venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	cupaloy.SnapshotT(t, rec.Body.String())
}
