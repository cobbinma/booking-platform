// +build !unit

package postgres_test

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/customer_api/internal/postgres"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/customer/api"
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

	var repository api.CustomerAPIServer
	var closeDB func(*zap.SugaredLogger)

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		p, c, err := postgres.
			NewPostgres(log, postgres.WithDatabaseURL(pgURL),
				postgres.WithMigrationsSourceURL("file://migrations"))
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

func suite(repository api.CustomerAPIServer) []test {
	return []test{
		{
			name: "is not administrator",
			test: func(t *testing.T) {
				resp, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					VenueId: "f982066f-1289-4317-83c1-d415dd4982c9",
					Email:   "test@test.com",
				})
				if err != nil {
					t.Fatalf("did not expect error from is admin, got '%s'", err)
				}

				if resp.IsAdmin != false {
					t.Errorf("expected is admin == false, got %v", resp.IsAdmin)
				}
			},
		},
		{
			name: "add administrator",
			test: func(t *testing.T) {
				venueID := "f982066f-1289-4317-83c1-d415dd4982c9"
				email := "test@test.com"
				resp, err := repository.AddAdmin(context.Background(), &api.AddAdminRequest{
					VenueId: venueID,
					Email:   email,
				})
				if err != nil {
					t.Fatalf("did not expect error from add admin, got '%s'", err)
				}

				if resp.VenueId != venueID {
					t.Errorf("expected is venueID == '%s', got '%s'", venueID, resp.VenueId)
				}

				if resp.Email != email {
					t.Errorf("expected is venueID == '%s', got '%s'", email, resp.Email)
				}
			},
		},
		{
			name: "is administrator",
			test: func(t *testing.T) {
				resp, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					VenueId: "f982066f-1289-4317-83c1-d415dd4982c9",
					Email:   "test@test.com",
				})
				if err != nil {
					t.Fatalf("did not expect error from is admin, got '%s'", err)
				}

				if resp.IsAdmin != true {
					t.Errorf("expected is admin == true, got %v", resp.IsAdmin)
				}
			},
		},
		{
			name: "remove administrator",
			test: func(t *testing.T) {
				venueID := "f982066f-1289-4317-83c1-d415dd4982c9"
				email := "test@test.com"
				resp, err := repository.RemoveAdmin(context.Background(), &api.RemoveAdminRequest{
					VenueId: venueID,
					Email:   email,
				})
				if err != nil {
					t.Fatalf("did not expect error from remove admin, got '%s'", err)
				}
				if resp.Email != email {
					t.Errorf("expected is venueID == '%s', got '%s'", email, resp.Email)
				}

				admin, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					VenueId: venueID,
					Email:   email,
				})
				if err != nil {
					t.Fatalf("did not expect error from is admin, got '%s'", err)
				}

				if admin.IsAdmin != false {
					t.Errorf("expected is admin == false, got %v", admin.IsAdmin)
				}
			},
		},
	}
}