// +build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking/lib/venue_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/venue_api/models"
	"github.com/cobbinma/booking/lib/venue_api/repositories/postgres"
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

func TestCreateVenue(t *testing.T) {
	name := "hop and vine"
	day := 1
	opens := "07:00"
	closes := "22:00"
	venueJSON := fmt.Sprintf(`{"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, name, day, opens, closes)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(venueJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	createVenue := handlers.CreateVenue(repository)
	err := createVenue(c)
	if err != nil {
		t.Errorf("error returned from create venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var vr venueResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &vr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	cupaloy.SnapshotT(t, vr)
}

func TestCreateGetVenue(t *testing.T) {
	name := "kings arms"
	day := 2
	opens := "08:00"
	closes := "21:00"
	venueJSON := fmt.Sprintf(`{"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, name, day, opens, closes)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(venueJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	createVenue := handlers.CreateVenue(repository)
	err := createVenue(c)
	if err != nil {
		t.Errorf("error returned from create venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var vr venueResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &vr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	id := vr.ID
	vr = venueResponse{}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	getVenue := handlers.GetVenue(repository)
	if err := getVenue(c); err != nil {
		t.Errorf("error returned from get venue handler : %s", err)
		return
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &vr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	cupaloy.SnapshotT(t, vr)
}

func TestCreateDeleteGetVenue(t *testing.T) {
	name := "case is altered"
	day := 1
	opens := "10:00"
	closes := "23:00"
	venueJSON := fmt.Sprintf(`{"name":"%s","openingHours":[{"dayOfWeek":%v,"opens":"%s","closes":"%s"}]}`, name, day, opens, closes)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(venueJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	createVenue := handlers.CreateVenue(repository)
	err := createVenue(c)
	if err != nil {
		t.Errorf("error returned from create venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusCreated)
		return
	}

	var vr venueResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &vr); err != nil {
		t.Errorf("could not unmarshall response : %s", err)
		return
	}

	id := vr.ID
	vr = venueResponse{}

	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/venues/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	deleteVenue := handlers.DeleteVenue(repository)
	if err := deleteVenue(c); err != nil {
		t.Errorf("error returned from get venue handler : %s", err)
		return
	}

	if rec.Code != http.StatusOK {
		t.Errorf("response code '%v' was not expected '%v'", rec.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
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
	}
}

type venueResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OpeningHours []struct {
		DayOfWeek int    `json:"dayOfWeek"`
		Opens     string `json:"opens"`
		Closes    string `json:"closes"`
	} `json:"openingHours"`
}
