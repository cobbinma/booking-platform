package postgres

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	"github.com/golang-migrate/migrate/v4"
	pgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/url"
	"os"
)

var _ Repository = (*client)(nil)

type Repository interface {
	api.TableAPIServer
	api.VenueAPIServer
}

type client struct {
	db               *sqlx.DB
	log              *zap.SugaredLogger
	pgURL            *url.URL
	migrationsSource string
}

func NewPostgres(log *zap.SugaredLogger, options ...func(*client)) (Repository, func(log *zap.SugaredLogger), error) {
	c := &client{log: log}
	for i := range options {
		options[i](c)
	}

	if c.pgURL == nil {
		u := os.Getenv("DATABASE_URL")
		if u == "" {
			return nil, nil, fmt.Errorf("environment variable 'DATABASE_URL' is not set")
		}
		p, err := url.Parse(u)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse 'DATABASE_URL'")
		}
		c.pgURL = p
	}

	if c.migrationsSource == "" {
		c.migrationsSource = "file://migrations"
	}

	db, err := sqlx.Connect("postgres", c.pgURL.String())
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to database : %w", err)
	}
	c.db = db

	if err := c.migrate(); err != nil {
		return nil, nil, fmt.Errorf("could not migrate : %w", err)
	}

	return c, func(log *zap.SugaredLogger) {
		if err := c.db.Close(); err != nil {
			c.log.Errorf("could not close database connection : %s", err)
		}
	}, nil
}

func (c client) GetTables(ctx context.Context, request *api.GetTablesRequest) (*api.GetTablesResponse, error) {
	panic("implement me")
}

func (c client) AddTable(ctx context.Context, request *api.AddTableRequest) (*models.Table, error) {
	panic("implement me")
}

func (c client) GetVenue(ctx context.Context, request *api.GetVenueRequest) (*models.Venue, error) {
	panic("implement me")
}

func (c client) CreateVenue(ctx context.Context, request *api.CreateVenueRequest) (*models.Venue, error) {
	panic("implement me")
}

func (c *client) migrate() error {
	driver, err := pgres.WithInstance(c.db.DB, &pgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver : %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("error instantiating migrate : %w", err)
	}

	version, dirty, _ := m.Version()
	c.log.Infof("database version %d, dirty %t", version, dirty)

	c.log.Infof("starting migration")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("an error occurred while syncing the database.. %w", err)
	}

	c.log.Infof("migration successfully")
	return nil
}
