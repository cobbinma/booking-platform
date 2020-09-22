// +build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/table_api/gateways/venueAPI"
	"github.com/cobbinma/booking/lib/table_api/models"
	"github.com/cobbinma/booking/lib/table_api/repositories/postgres"
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

func TestCreateTable(t *testing.T) {
	name := "2"
	capacity := 4
	venueID := 1
	venueName := "hop and vine"
	day := 1
	opens := "07:00"
	closes := "22:00"
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, day, opens, closes)
	tableJSON := fmt.Sprintf(`{"name": "%v","capacity": %v}`, name, capacity)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tableJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer srv.Close()

	createTable := handlers.VenueMiddleware(handlers.CreateTable(repository), venueAPI.NewVenueAPI(srv.URL))

	err := createTable(c)
	if err != nil {
		t.Errorf("error returned from create table handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var tr models.Table
	if err := json.Unmarshal(rec.Body.Bytes(), &tr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	if tr.ID == 0 {
		t.Errorf("table ID should not equal 0")
	}

	if tr.Name != name {
		t.Errorf("name was '%s', expected '%s'", tr.Name, name)
	}

	if tr.Capacity != capacity {
		t.Errorf("capacity was '%v', expected '%v'", tr.Capacity, capacity)
	}
}

func TestCreateGetTable(t *testing.T) {
	name := "2"
	capacity := 4
	venueID := 1
	venueName := "hop and vine"
	day := 1
	opens := "07:00"
	closes := "22:00"
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, day, opens, closes)
	tableJSON := fmt.Sprintf(`{"name": "%v","capacity": %v}`, name, capacity)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tableJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer srv.Close()

	createTable := handlers.VenueMiddleware(handlers.CreateTable(repository), venueAPI.NewVenueAPI(srv.URL))

	err := createTable(c)
	if err != nil {
		t.Errorf("error returned from create table handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var tr models.Table
	if err := json.Unmarshal(rec.Body.Bytes(), &tr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	tableID := tr.ID
	tr = models.Table{}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables/:id")
	c.SetParamNames("venue_id", "id")
	c.SetParamValues(strconv.Itoa(venueID), tableID.String())

	getTable := handlers.VenueMiddleware(handlers.GetTable(repository), venueAPI.NewVenueAPI(srv.URL))

	if err := getTable(c); err != nil {
		t.Errorf("error returned from get table handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &tr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	checkTable(t, tr, tableID, name, capacity)
}

func TestCreateGetTables(t *testing.T) {
	name := "6"
	capacity := 4
	venueID := 1
	venueName := "hop and vine"
	day := 1
	opens := "07:00"
	closes := "22:00"
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, day, opens, closes)
	tableJSON := fmt.Sprintf(`{"name": "%v","capacity": %v}`, name, capacity)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tableJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer srv.Close()

	createTable := handlers.VenueMiddleware(handlers.CreateTable(repository), venueAPI.NewVenueAPI(srv.URL))

	err := createTable(c)
	if err != nil {
		t.Errorf("error returned from create table handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var tr models.Table
	if err := json.Unmarshal(rec.Body.Bytes(), &tr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	tableID := tr.ID
	tr = models.Table{}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables/:id")
	c.SetParamNames("venue_id", "id")
	c.SetParamValues(strconv.Itoa(venueID), tableID.String())

	getTables := handlers.VenueMiddleware(handlers.GetTables(repository), venueAPI.NewVenueAPI(srv.URL))

	if err := getTables(c); err != nil {
		t.Errorf("error returned from get tables handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var tables []models.Table
	if err := json.Unmarshal(rec.Body.Bytes(), &tables); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	for i := range tables {
		if tables[i].ID == tableID {
			checkTable(t, tables[i], tableID, name, capacity)
			return
		}
	}
	t.Errorf("could not find table in get tables response")
}

func TestCreateGetTablesWithCapacity(t *testing.T) {
	name := "17"
	capacity := 4
	venueID := 1
	venueName := "hop and vine"
	day := 1
	opens := "07:00"
	closes := "22:00"
	venueJSON := fmt.Sprintf(`{"id":%v,"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, venueID, venueName, day, opens, closes)
	tableJSON := fmt.Sprintf(`{"name": "%v","capacity": %v}`, name, capacity)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tableJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables")
	c.SetParamNames("venue_id")
	c.SetParamValues(strconv.Itoa(venueID))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(venueJSON))
	}))
	defer srv.Close()

	createTable := handlers.VenueMiddleware(handlers.CreateTable(repository), venueAPI.NewVenueAPI(srv.URL))

	err := createTable(c)
	if err != nil {
		t.Errorf("error returned from create table handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var tr models.Table
	if err := json.Unmarshal(rec.Body.Bytes(), &tr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	tableID := tr.ID
	tr = models.Table{}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:venue_id/tables/capacity/:amount")
	c.SetParamNames("venue_id", "amount")
	c.SetParamValues(strconv.Itoa(venueID), strconv.Itoa(capacity))

	getTables := handlers.VenueMiddleware(handlers.GetTablesWithCapacity(repository), venueAPI.NewVenueAPI(srv.URL))

	if err := getTables(c); err != nil {
		t.Errorf("error returned from get tables handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var tables []models.Table
	if err := json.Unmarshal(rec.Body.Bytes(), &tables); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	for i := range tables {
		if tables[i].ID == tableID {
			checkTable(t, tables[i], tableID, name, capacity)
			return
		}
	}
	t.Errorf("could not find table in get tables response")
}

func checkTable(t *testing.T, tr models.Table, tableID models.TableID, name string, capacity int) {
	if tr.ID != tableID {
		t.Errorf("id was '%v', expected '%v'", tr.ID, tableID)
	}

	if tr.Name != name {
		t.Errorf("name was '%s', expected '%s'", tr.Name, name)
	}

	if tr.Capacity != capacity {
		t.Errorf("capacity was '%v', expected '%v'", tr.Capacity, capacity)
	}
}
