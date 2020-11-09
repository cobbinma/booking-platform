package handlers_test

import (
	"github.com/cobbinma/booking/lib/venue_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/venue_api/repositories/fakeRepository"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestVenueNotFound(t *testing.T) {
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
}
