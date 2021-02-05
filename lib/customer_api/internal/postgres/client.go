package postgres

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/customer/api"
	"github.com/golang-migrate/migrate/v4"
	pgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"os"
)

const (
	AdminsTable = "admins"
)

var _ api.CustomerAPIServer = (*client)(nil)

type client struct {
	db               *sqlx.DB
	log              *zap.SugaredLogger
	pgURL            *url.URL
	migrationsSource string
}

func NewPostgres(log *zap.SugaredLogger, options ...func(*client)) (api.CustomerAPIServer, func(log *zap.SugaredLogger), error) {
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
		c.migrationsSource = "file://internal/postgres/migrations"
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

func (c client) IsAdmin(ctx context.Context, req *api.IsAdminRequest) (*api.IsAdminResponse, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From(AdminsTable).
		Where(sq.And{sq.Eq{"venue_id": req.VenueId}, sq.Eq{"email": req.Email}}).ToSql()
	if err != nil {
		c.log.Errorw("could not construct sql", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	row := c.db.QueryRow(sql, args...)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql2.ErrNoRows) {
			return &api.IsAdminResponse{IsAdmin: false}, nil
		}

		c.log.Errorw("could not query row", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	return &api.IsAdminResponse{IsAdmin: true}, nil
}

func (c client) AddAdmin(ctx context.Context, req *api.AddAdminRequest) (*api.AddAdminResponse, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(AdminsTable).Columns("id", "venue_id", "email").
		Values(uuid.New().String(), req.VenueId, req.Email).ToSql()
	if err != nil {
		c.log.Errorw("could not construct sql", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	_, err = c.db.Exec(sql, args...)
	if err != nil {
		c.log.Errorw("could not insert row", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not insert row")
	}

	return &api.AddAdminResponse{
		VenueId: req.VenueId,
		Email:   req.Email,
	}, nil
}

func (c client) RemoveAdmin(ctx context.Context, req *api.RemoveAdminRequest) (*api.RemoveAdminResponse, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(AdminsTable).
		Where(sq.And{sq.Eq{"venue_id": req.VenueId}, sq.Eq{"email": req.Email}}).ToSql()
	if err != nil {
		c.log.Errorw("could not construct sql", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	_, err = c.db.Exec(sql, args...)
	if err != nil {
		c.log.Errorw("could not delete row", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not delete row")
	}

	return &api.RemoveAdminResponse{Email: req.Email}, nil
}

func (c *client) migrate() error {
	driver, err := pgres.WithInstance(c.db.DB, &pgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver : %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		c.migrationsSource,
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
