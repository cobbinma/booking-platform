package handlers_test

import (
	"fmt"
	"github.com/cobbinma/booking/lib/venue_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/venue_api/repositories/fakeRepository"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidOpeningHours(t *testing.T) {
	repository := fakeRepository.NewFakeRepository()

	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		fmt.Sprintf(`{"name":"hop and vine","openingHours":[{"dayOfWeek":1,"opens":15,"closes":"22:00"}]}`)))
	c := e.NewContext(req, rec)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	createVenue := handlers.CreateVenue(repository)

	if err := createVenue(c); err != nil {
		t.Errorf("error returned from create venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusBadRequest)
		return
	}
}
