// +build !unit

package postgres_test

import (
	"context"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"github.com/cobbinma/booking-platform/lib/venue_api/internal/postgres"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"
	"net"
	"net/url"
	"runtime"
	"testing"
	"time"
)

func Test_Repository(t *testing.T) {
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
		t.Errorf("could not connect to docker : %s", err)
		return
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
		t.Errorf("Could start postgres container : %s", err)
		return
	}
	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			t.Errorf("Could not purge resource : %s", err)
		}
	}()

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	logger, err := zap.NewProduction()
	if err != nil {
		t.Errorf("Could start construct zap logger : %s", err)
		return
	}
	log := logger.Sugar()
	defer func(log *zap.SugaredLogger) {
		if err := logger.Sync(); err != nil {
			log.Errorf("could not sync logger : %s", err)
		}
	}(log)

	var repository postgres.Repository
	var closeDB func(*zap.SugaredLogger)

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		p, c, err := postgres.
			NewPostgres(log, postgres.WithDatabaseURL(pgURL),
				postgres.WithMigrationsSourceURL("file://migrations"),
				postgres.WithStaticUUIDGenerator("b31a9f99-3f64-4ee9-af27-45b2acd36d86"))
		if err != nil {
			return err
		}
		closeDB = c
		repository = p
		return nil
	})
	if err != nil {
		t.Errorf("could not connect to postgres server : %s", err)
		return
	}
	defer closeDB(log)

	s := suite(repository)
	for i := range s {
		t.Run(s[i].name, s[i].test)
	}
}

type test struct {
	name string
	test func(t *testing.T)
}

func suite(repository postgres.Repository) []test {
	return []test{
		{
			name: "add venue successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				venue, err := repository.CreateVenue(ctx, &api.CreateVenueRequest{
					Name: "Test Venue",
					OpeningHours: []*models.OpeningHoursSpecification{
						{
							DayOfWeek:    1,
							Opens:        "10:00",
							Closes:       "20:00",
							ValidFrom:    "",
							ValidThrough: "",
						},
						{
							DayOfWeek:    2,
							Opens:        "10:00",
							Closes:       "22:00",
							ValidFrom:    "",
							ValidThrough: "",
						},
					},
				})
				if err != nil {
					t.Fatalf("did not expect error, got '%s'", err)
				}

				cupaloy.SnapshotT(t, venue)
			},
		},
		{
			name: "get venue successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				venues, err := repository.GetVenue(ctx, &api.GetVenueRequest{Id: "b31a9f99-3f64-4ee9-af27-45b2acd36d86"})
				if err != nil {
					t.Fatalf("did not expect error, got '%s'", err)
				}

				cupaloy.SnapshotT(t, venues)
			},
		},
		{
			name: "add table successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				table, err := repository.AddTable(ctx, &api.AddTableRequest{
					VenueId: "1e182275-00e4-4334-b765-5cca09d5e548",
					Table: &models.Table{
						Id:       "66d499a2-75e2-400c-a9aa-43f6d08c5d2b",
						Name:     "test table",
						Capacity: 4,
					},
				})
				if err != nil {
					t.Fatalf("did not expect error, got '%s'", err)
				}

				cupaloy.SnapshotT(t, table)
			},
		},
		{
			name: "get tables successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				table, err := repository.GetTables(ctx, &api.GetTablesRequest{
					VenueId: "1e182275-00e4-4334-b765-5cca09d5e548"})
				if err != nil {
					t.Fatalf("did not expect error, got '%s'", err)
				}

				cupaloy.SnapshotT(t, table)
			},
		},
	}
}
