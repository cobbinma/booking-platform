// +build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/booking_api/gateways/tableAPI"
	"github.com/cobbinma/booking/lib/booking_api/gateways/venueAPI"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/cobbinma/booking/lib/booking_api/repositories/postgres"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	repository models.Repository
	now        = time.Now()
	name       = "2"
	venueID    = 1
	venueName  = "hop and vine"
	opens      = "07:00"
	closes     = "22:00"
	customer   = "example@example.com"
	people     = 2
	tableId    = 1
)

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceColors:     true,
	}

	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("myuser", "mypass"),
		Path:   "mydatabase",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Fatal("Could not connect to docker")
	}

	pw, _ := pgURL.User.Password()
	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + pgURL.User.Username(),
			"POSTGRES_PASSWORD=" + pw,
			"POSTGRES_DB=" + pgURL.Path,
		},
	}

	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		log.WithError(err).Fatal("Could start postgres container")
	}
	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.WithError(err).Error("Could not purge resource")
		}
	}()

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: log.Writer(),
		ErrorStream:  log.Writer(),
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.WithError(err).Fatal("Could not connect to postgres container log output")
	}
	defer func() {
		err = logWaiter.Close()
		if err != nil {
			log.WithError(err).Error("Could not close container log")
		}
		err = logWaiter.Wait()
		if err != nil {
			log.WithError(err).Error("Could not wait for container log to close")
		}
	}()

	pool.MaxWait = 10 * time.Second

	var closeDB func() error
	var dbClient postgres.DBClient

	err = pool.Retry(func() error {
		dbClient, closeDB, err = postgres.NewDBClient(pgURL)
		if err != nil {
			return err
		}
		return dbClient.DB().Ping()
	})
	if err != nil {
		log.WithError(err).Fatal("Could not connect to postgres server")
	}
	defer func() {
		if err := closeDB(); err != nil {
			log.Error("could not close database : ", err)
		}
	}()

	repository = postgres.NewPostgres(dbClient)
	if err := repository.Migrate(context.Background(), "file://../migrations"); err != nil {
		log.Fatal("could not migrate : ", err)
	}

	code = m.Run()
}

func TestBookingQuery(t *testing.T) {
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	tomorrowDate := models.Date(tomorrow)
	tomorrowStr := tomorrowDate.Format(models.DateFormat)
	startsAt := time.Date(now.Year(), now.Month(), now.Day()+1, 18, 0, 0, 0, time.UTC)
	startsAtStr := startsAt.Format("2006-01-02T15:04:05Z")
	endsAt := time.Date(now.Year(), now.Month(), now.Day()+1, 20, 0, 0, 0, time.UTC)
	endsAtStr := endsAt.Format("2006-01-02T15:04:05Z")
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, int(startsAt.Weekday()), opens, closes)
	queryJSON := fmt.Sprintf(`{"customer_id":"%s","people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, people, tomorrowStr, startsAtStr, endsAtStr)
	tableJSON := fmt.Sprintf(`[{"id":%v,"name":"%v","capacity":%v}]`, tableId, name, people)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(queryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/slot")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	venueSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer venueSrv.Close()

	tableSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))
	defer tableSrv.Close()

	queryBooking := handlers.VenueMiddleware(handlers.BookingQuery(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err := queryBooking(c)
	if err != nil {
		t.Errorf("error returned from booking query handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	var nb models.Slot
	if err := json.Unmarshal(rec.Body.Bytes(), &nb); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	if nb.CustomerID != models.CustomerID(customer) {
		t.Errorf("customer was '%s', expected '%s'", nb.CustomerID, customer)
	}

	if nb.TableID != models.TableID(tableId) {
		t.Errorf("table id was '%v', expected '%v'", nb.TableID, tableId)
	}

	if nb.Date != tomorrowDate {
		t.Errorf("date was '%v', expected '%v'", nb.Date, tomorrowDate)
	}

	if nb.StartsAt != startsAt {
		t.Errorf("starts at was '%v', expected '%v'", nb.StartsAt, startsAt)
	}

	if nb.EndsAt != endsAt {
		t.Errorf("ends at was '%v', expected '%v'", nb.EndsAt, endsAt)
	}
}

func TestBookingQueryCreateBooking(t *testing.T) {
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	tomorrowDate := models.Date(tomorrow)
	tomorrowStr := tomorrowDate.Format(models.DateFormat)
	startsAt := time.Date(now.Year(), now.Month(), now.Day()+1, 18, 0, 0, 0, time.UTC)
	startsAtStr := startsAt.Format("2006-01-02T15:04:05Z")
	endsAt := time.Date(now.Year(), now.Month(), now.Day()+1, 20, 0, 0, 0, time.UTC)
	endsAtStr := endsAt.Format("2006-01-02T15:04:05Z")
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, int(startsAt.Weekday()), opens, closes)
	queryJSON := fmt.Sprintf(`{"customer_id":"%s","people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, people, tomorrowStr, startsAtStr, endsAtStr)
	tableJSON := fmt.Sprintf(`[{"id":%v,"name":"%v","capacity":%v}]`, tableId, name, people)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(queryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/slot")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	venueSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer venueSrv.Close()

	tableSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))
	defer tableSrv.Close()

	queryBooking := handlers.VenueMiddleware(handlers.BookingQuery(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err := queryBooking(c)
	if err != nil {
		t.Errorf("error returned from booking query handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	var nb models.Slot
	if err := json.Unmarshal(rec.Body.Bytes(), &nb); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	bookingJSON := fmt.Sprintf(`{"customer_id":"%s","table_id":%v,"people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, nb.TableID, people, tomorrowStr, startsAtStr, endsAtStr)

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bookingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/booking")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	tableJSON = fmt.Sprintf(`{"id":%v,"name":"%v","capacity":%v}`, tableId, name, people)
	tableSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))

	createBooking := handlers.VenueMiddleware(handlers.CreateBooking(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err = createBooking(c)
	if err != nil {
		t.Errorf("error returned from create booking handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var booking models.Booking
	if err := json.Unmarshal(rec.Body.Bytes(), &booking); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}
}

func TestBookingQueryCreateBookingGetBookingByDate(t *testing.T) {
	twoDays := time.Date(now.Year(), now.Month(), now.Day()+2, 0, 0, 0, 0, time.UTC)
	twoDaysStr := twoDays.Format(models.DateFormat)
	twoDaysDate := models.Date(twoDays)
	startsAt := time.Date(now.Year(), now.Month(), now.Day()+2, 18, 0, 0, 0, time.UTC)
	startsAtStr := startsAt.Format("2006-01-02T15:04:05Z")
	endsAt := time.Date(now.Year(), now.Month(), now.Day()+2, 20, 0, 0, 0, time.UTC)
	endsAtStr := endsAt.Format("2006-01-02T15:04:05Z")
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, int(startsAt.Weekday()), opens, closes)
	queryJSON := fmt.Sprintf(`{"customer_id":"%s","people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, people, twoDaysStr, startsAtStr, endsAtStr)
	tableJSON := fmt.Sprintf(`[{"id":%v,"name":"%v","capacity":%v}]`, tableId, name, people)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(queryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/slot")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	venueSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer venueSrv.Close()

	tableSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))
	defer tableSrv.Close()

	queryBooking := handlers.VenueMiddleware(handlers.BookingQuery(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err := queryBooking(c)
	if err != nil {
		t.Errorf("error returned from booking query handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	var nb models.Slot
	if err := json.Unmarshal(rec.Body.Bytes(), &nb); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	bookingJSON := fmt.Sprintf(`{"customer_id":"%s","table_id":%v,"people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, nb.TableID, people, twoDaysStr, startsAtStr, endsAtStr)

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bookingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/booking")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	tableJSON = fmt.Sprintf(`{"id":%v,"name":"%v","capacity":%v}`, tableId, name, people)
	tableSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))

	createBooking := handlers.VenueMiddleware(handlers.CreateBooking(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err = createBooking(c)
	if err != nil {
		t.Errorf("error returned from create booking handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bookingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/booking/date/:date")
	c.SetParamNames("venue_id", "date")
	c.SetParamValues(strconv.Itoa(venueID), twoDaysStr)

	getBooking := handlers.VenueMiddleware(handlers.GetBookingsByDate(repository), venueAPI.NewVenueAPI(venueSrv.URL))

	err = getBooking(c)
	if err != nil {
		t.Errorf("error returned from get booking handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	var bookings []models.Booking
	if err := json.Unmarshal(rec.Body.Bytes(), &bookings); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	if bookings[0].CustomerID != models.CustomerID(customer) {
		t.Errorf("customer was '%s', expected '%s'", bookings[0].CustomerID, customer)
	}

	if bookings[0].TableID != models.TableID(tableId) {
		t.Errorf("table id was '%v', expected '%v'", bookings[0].TableID, tableId)
	}

	if bookings[0].Date != twoDaysDate {
		t.Errorf("date was '%v', expected '%v'", bookings[0].Date, twoDaysStr)
	}

	if bookings[0].StartsAt != startsAt {
		t.Errorf("starts at was '%v', expected '%v'", bookings[0].StartsAt, startsAt)
	}

	if bookings[0].EndsAt != endsAt {
		t.Errorf("ends at was '%v', expected '%v'", bookings[0].EndsAt, endsAt)
	}
}

func TestBookingQueryCreateBookingDeleteBookingGetBookingByDate(t *testing.T) {
	threeDays := time.Date(now.Year(), now.Month(), now.Day()+3, 0, 0, 0, 0, time.UTC)
	threeDaysStr := threeDays.Format(models.DateFormat)
	startsAt := time.Date(now.Year(), now.Month(), now.Day()+3, 18, 0, 0, 0, time.UTC)
	startsAtStr := startsAt.Format("2006-01-02T15:04:05Z")
	endsAt := time.Date(now.Year(), now.Month(), now.Day()+3, 20, 0, 0, 0, time.UTC)
	endsAtStr := endsAt.Format("2006-01-02T15:04:05Z")
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, int(startsAt.Weekday()), opens, closes)
	queryJSON := fmt.Sprintf(`{"customer_id":"%s","people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, people, threeDaysStr, startsAtStr, endsAtStr)
	tableJSON := fmt.Sprintf(`[{"id":%v,"name":"%v","capacity":%v}]`, tableId, name, people)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(queryJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/slot")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	venueSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer venueSrv.Close()

	tableSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))
	defer tableSrv.Close()

	queryBooking := handlers.VenueMiddleware(handlers.BookingQuery(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err := queryBooking(c)
	if err != nil {
		t.Errorf("error returned from booking query handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	var nb models.Slot
	if err := json.Unmarshal(rec.Body.Bytes(), &nb); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	bookingJSON := fmt.Sprintf(`{"customer_id":"%s","table_id":%v,"people":%v,"date":"%s","starts_at":"%s","ends_at":"%s"}`, customer, nb.TableID, people, threeDaysStr, startsAtStr, endsAtStr)

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bookingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/booking")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	tableJSON = fmt.Sprintf(`{"id":%v,"name":"%v","capacity":%v}`, tableId, name, people)
	tableSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tableJSON))
	}))

	createBooking := handlers.VenueMiddleware(handlers.CreateBooking(repository, tableAPI.NewTableAPI(tableSrv.URL)), venueAPI.NewVenueAPI(venueSrv.URL))

	err = createBooking(c)
	if err != nil {
		t.Errorf("error returned from create booking handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var booking models.Booking
	if err := json.Unmarshal(rec.Body.Bytes(), &booking); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	req = httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(bookingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/booking/:id")
	c.SetParamNames("venue_id", "id")
	c.SetParamValues(strconv.Itoa(venueID), strconv.Itoa(booking.ID))

	deleteBooking := handlers.VenueMiddleware(handlers.DeleteBooking(repository), venueAPI.NewVenueAPI(venueSrv.URL))

	err = deleteBooking(c)
	if err != nil {
		t.Errorf("error returned from get booking handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	req = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(bookingJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/booking/date/:date")
	c.SetParamNames("venue_id", "date")
	c.SetParamValues(strconv.Itoa(venueID), threeDaysStr)

	getBooking := handlers.VenueMiddleware(handlers.GetBookingsByDate(repository), venueAPI.NewVenueAPI(venueSrv.URL))

	err = getBooking(c)
	if err != nil {
		t.Errorf("error returned from get booking handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
		return
	}

	var bookings []models.Booking
	if err := json.Unmarshal(rec.Body.Bytes(), &bookings); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	if len(bookings) != 0 {
		t.Errorf("should not have found bookings")
		return
	}
}
