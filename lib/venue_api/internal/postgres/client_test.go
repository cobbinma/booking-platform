// +build !unit

package postgres_test

import (
	"context"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"github.com/cobbinma/booking-platform/lib/venue_api/internal/postgres"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	require.NoError(t, err)

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
	require.NoError(t, err)
	defer func() {
		require.NoError(t, pool.Purge(resource))
	}()

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	log := zap.NewNop().Sugar()

	var repository api.VenueAPIServer
	var closeDB func(*zap.SugaredLogger)

	pool.MaxWait = 10 * time.Second
	require.NoError(t, pool.Retry(func() error {
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
	}))
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

func suite(repository api.VenueAPIServer) []test {
	const UUID = "b31a9f99-3f64-4ee9-af27-45b2acd36d86"
	const Slug = "test-venue"
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
					Slug: Slug,
				})
				require.NoError(t, err)

				cupaloy.SnapshotT(t, venue)
			},
		},
		{
			name: "get venue successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				venues, err := repository.GetVenue(ctx, &api.GetVenueRequest{Id: UUID})
				require.NoError(t, err)

				cupaloy.SnapshotT(t, venues)
			},
		},
		{
			name: "get venue by slug successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				venues, err := repository.GetVenue(ctx, &api.GetVenueRequest{Slug: Slug})
				require.NoError(t, err)

				cupaloy.SnapshotT(t, venues)
			},
		},
		{
			name: "update venue opening hours",
			test: func(t *testing.T) {
				ctx := context.Background()
				venues, err := repository.UpdateOpeningHours(ctx, &api.UpdateOpeningHoursRequest{VenueId: UUID, OpeningHours: []*models.OpeningHoursSpecification{
					{
						DayOfWeek:    2,
						Opens:        "11:00",
						Closes:       "22:00",
						ValidFrom:    "",
						ValidThrough: "",
					},
					{
						DayOfWeek:    3,
						Opens:        "10:30",
						Closes:       "23:00",
						ValidFrom:    "",
						ValidThrough: "",
					},
				}})
				require.NoError(t, err)

				cupaloy.SnapshotT(t, venues)
			},
		},
		{
			name: "add table successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				table, err := repository.AddTable(ctx, &api.AddTableRequest{
					VenueId:  UUID,
					Name:     "test table",
					Capacity: 4,
				})
				require.NoError(t, err)

				cupaloy.SnapshotT(t, table)
			},
		},
		{
			name: "get tables successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				table, err := repository.GetTables(ctx, &api.GetTablesRequest{
					VenueId: UUID})
				require.NoError(t, err)

				cupaloy.SnapshotT(t, table)
			},
		},
		{
			name: "remove not found table",
			test: func(t *testing.T) {
				ctx := context.Background()
				_, err := repository.RemoveTable(ctx, &api.RemoveTableRequest{
					VenueId: UUID,
					TableId: uuid.New().String(),
				})
				assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "remove table successfully",
			test: func(t *testing.T) {
				ctx := context.Background()
				removed, err := repository.RemoveTable(ctx, &api.RemoveTableRequest{
					VenueId: UUID,
					TableId: UUID,
				})
				require.NoError(t, err)

				assert.Equal(t, &models.Table{
					Id:       UUID,
					Name:     "test table",
					Capacity: 4,
				}, removed)
			},
		},
		{
			name: "is not administrator",
			test: func(t *testing.T) {
				resp, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					VenueId: UUID,
					Email:   "test@test.com",
				})
				require.NoError(t, err)

				assert.Equal(t, false, resp.IsAdmin)
			},
		},
		{
			name: "add administrator",
			test: func(t *testing.T) {
				venueID := UUID
				email := "test@test.com"
				resp, err := repository.AddAdmin(context.Background(), &api.AddAdminRequest{
					VenueId: venueID,
					Email:   email,
				})
				require.NoError(t, err)

				assert.Equal(t, venueID, resp.VenueId)
				assert.Equal(t, email, resp.Email)
			},
		},
		{
			name: "is administrator",
			test: func(t *testing.T) {
				resp, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					VenueId: UUID,
					Email:   "test@test.com",
				})
				require.NoError(t, err)

				assert.Equal(t, true, resp.IsAdmin)
			},
		},
		{
			name: "is administrator by slug",
			test: func(t *testing.T) {
				resp, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					Slug:  Slug,
					Email: "test@test.com",
				})
				require.NoError(t, err)

				assert.Equal(t, true, resp.IsAdmin)
			},
		},
		{
			name: "get administrators",
			test: func(t *testing.T) {
				resp, err := repository.GetAdmins(context.Background(), &api.GetAdminsRequest{VenueId: UUID})
				require.NoError(t, err)

				require.Equal(t, 1, len(resp.Admins))
				assert.Equal(t, "test@test.com", resp.Admins[0])
			},
		},
		{
			name: "remove administrator",
			test: func(t *testing.T) {
				venueID := UUID
				email := "test@test.com"
				resp, err := repository.RemoveAdmin(context.Background(), &api.RemoveAdminRequest{
					VenueId: venueID,
					Email:   email,
				})
				require.NoError(t, err)
				require.Equal(t, email, resp.Email)

				admin, err := repository.IsAdmin(context.Background(), &api.IsAdminRequest{
					VenueId: venueID,
					Email:   email,
				})
				require.NoError(t, err)

				assert.Equal(t, false, admin.IsAdmin)
			},
		},
		{
			name: "get administrators none",
			test: func(t *testing.T) {
				resp, err := repository.GetAdmins(context.Background(), &api.GetAdminsRequest{VenueId: UUID})
				require.NoError(t, err)

				require.Equal(t, 0, len(resp.Admins))
			},
		},
	}
}
