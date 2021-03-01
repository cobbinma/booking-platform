package postgres

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
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
	OpeningHoursTable = "opening_hours"
	VenuesTable       = "venues"
	TablesTable       = "tables"
	AdminsTable       = "admins"
)

var _ api.VenueAPIServer = (*client)(nil)

type client struct {
	db               *sqlx.DB
	log              *zap.SugaredLogger
	pgURL            *url.URL
	migrationsSource string
	uuid             uuidGenerator
}

func NewPostgres(log *zap.SugaredLogger, options ...func(*client)) (api.VenueAPIServer, func(log *zap.SugaredLogger), error) {
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

	if c.uuid == nil {
		c.uuid = newRandomUUID()
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

func (c client) GetTables(ctx context.Context, req *api.GetTablesRequest) (*api.GetTablesResponse, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name", "capacity").
		From(TablesTable).Where(sq.Eq{"venue_id": req.VenueId}).ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build tables sql : %s", err)
	}

	tables := []*models.Table{}
	rows, err := c.db.Query(sql, args...)
	if err != nil && !errors.Is(err, sql2.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "could not query tables : %s", err)
	}
	if rows != nil {
		for rows.Next() {
			var capacity uint32
			var id, name string
			if err := rows.Scan(&id, &name, &capacity); err != nil {
				return nil, status.Errorf(codes.Internal, "could not scan tables row : %s", err)
			}
			tables = append(tables, &models.Table{
				Id:       id,
				Name:     name,
				Capacity: capacity,
			})
		}

		if err := rows.Err(); err != nil {
			return nil, status.Errorf(codes.Internal, "tables rows error : %s", err)
		}
	}

	return &api.GetTablesResponse{Tables: tables}, nil
}

func (c client) AddTable(ctx context.Context, req *api.AddTableRequest) (*models.Table, error) {
	id := c.uuid.UUID()
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(TablesTable).
		Columns("id", "venue_id", "name", "capacity").
		Values(id, req.VenueId, req.Name, req.Capacity).ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build table sql : %s", err)
	}

	if _, err := c.db.Exec(sql, args...); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert table : %s", err)
	}

	return &models.Table{
		Id:       id,
		Name:     req.Name,
		Capacity: req.Capacity,
	}, nil
}

func (c client) RemoveTable(ctx context.Context, req *api.RemoveTableRequest) (*models.Table, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name", "capacity").
		From(TablesTable).
		Where(sq.And{sq.Eq{"id": req.TableId}, sq.Eq{"venue_id": req.VenueId}}).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build select table sql : %s", err)
	}

	row := c.db.QueryRow(sql, args...)
	if row.Err() != nil {
		c.log.Errorw("could not query row", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	var id, name string
	var capacity uint32
	if err := c.db.QueryRow(sql, args...).Scan(&id, &name, &capacity); err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "could not find venue")
		}

		return nil, status.Errorf(codes.Internal, "could get find venue : %s", err)
	}

	sql, args, err = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(TablesTable).
		Where(sq.And{sq.Eq{"id": req.TableId}, sq.Eq{"venue_id": req.VenueId}}).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build select table sql : %s", err)
	}

	if _, err := c.db.Exec(sql, args...); err != nil {
		return nil, status.Errorf(codes.Internal, "could not delete table : %s", err)
	}

	return &models.Table{Id: id, Name: name, Capacity: capacity}, nil
}

func (c client) GetVenue(ctx context.Context, req *api.GetVenueRequest) (*models.Venue, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name").From(VenuesTable).
		Where(sq.Eq{"id": req.Id}).ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not venue build sql : %s", err)
	}

	var id, name string
	if err := c.db.QueryRow(sql, args...).Scan(&id, &name); err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "could not find venue")
		}

		return nil, status.Errorf(codes.Internal, "could get find venue : %s", err)
	}

	sql, args, err = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("day_of_week", "opens", "closes").
		From(OpeningHoursTable).Where(sq.Eq{"venue_id": req.Id}).ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build opening hours sql : %s", err)
	}

	hours := []*models.OpeningHoursSpecification{}
	rows, err := c.db.Query(sql, args...)
	if err != nil && !errors.Is(err, sql2.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "could not query opening hours : %s", err)
	}
	if rows != nil {
		for rows.Next() {
			var day_of_week uint32
			var opens, closes string
			if err := rows.Scan(&day_of_week, &opens, &closes); err != nil {
				return nil, status.Errorf(codes.Internal, "could not scan opening hours row : %s", err)
			}
			hours = append(hours, &models.OpeningHoursSpecification{
				DayOfWeek: day_of_week,
				Opens:     opens,
				Closes:    closes,
			})
		}

		if err := rows.Err(); err != nil {
			return nil, status.Errorf(codes.Internal, "opening hours rows error : %s", err)
		}
	}

	return &models.Venue{
		Id:           id,
		Name:         name,
		OpeningHours: hours,
	}, nil
}

func (c client) CreateVenue(ctx context.Context, req *api.CreateVenueRequest) (*models.Venue, error) {
	id := c.uuid.UUID()
	tx, err := c.db.Beginx()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not begin transaction : %s", err)
	}

	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(VenuesTable).
		Columns("id", "name").Values(id, req.GetName()).ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build venue sql : %s", err)
	}

	if _, err := tx.Exec(sql, args...); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert venue : %s", err)
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(OpeningHoursTable).
		Columns("venue_id", "day_of_week", "opens", "closes")

	for _, hours := range req.OpeningHours {
		builder = builder.Values(
			id,
			hours.DayOfWeek,
			hours.Opens,
			hours.Closes,
		)
	}

	sql, args, err = builder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not build opening_hours sql : %s", err)
	}

	if _, err := tx.Exec(sql, args...); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert opening hours : %s", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "could not commit transaction : %s", err)
	}

	return &models.Venue{
		Id:                  id,
		Name:                req.Name,
		OpeningHours:        req.OpeningHours,
		SpecialOpeningHours: nil,
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
		c.log.Errorw("could not query row", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	var count int
	if err := row.Scan(&count); err != nil {
		c.log.Errorw("could not scan row", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal database error")
	}

	if count != 1 {
		return &api.IsAdminResponse{IsAdmin: false}, nil
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
