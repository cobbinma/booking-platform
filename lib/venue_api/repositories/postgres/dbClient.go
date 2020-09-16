package postgres

import (
	"database/sql"
	"fmt"
	"github.com/cobbinma/booking/lib/venue_api/config"
	"github.com/cobbinma/booking/lib/venue_api/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBClient interface {
	DB() *sql.DB
	Exec(query string, args ...interface{}) (sql.Result, error)
	GetVenue(query string, args ...interface{}) (*models.Venue, error)
	GetOpeningHours(query string, args ...interface{}) ([]models.OpeningHours, error)
	BeginX() (*sqlx.Tx, error)
}

type dbClient struct {
	db *sqlx.DB
}

func NewDBClient() (*dbClient, func() error, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		config.DBHost,
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBSSLMode)

	driver := "postgres"

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open database : %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	dbc := &dbClient{db: db}

	return dbc, dbc.Close, nil
}

func (dbc *dbClient) Close() error {
	return dbc.DB().Close()
}

func (dbc *dbClient) DB() *sql.DB {
	return dbc.db.DB
}

func (dbc *dbClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dbc.db.Exec(query, args...)
}

func (dbc *dbClient) GetVenue(query string, args ...interface{}) (*models.Venue, error) {
	venue := &models.Venue{}
	if err := dbc.db.Get(venue, query, args...); err != nil {
		return nil, err
	}

	return venue, nil
}

func (dbc *dbClient) GetOpeningHours(query string, args ...interface{}) ([]models.OpeningHours, error) {
	oh := []models.OpeningHours{}
	if err := dbc.db.Select(&oh, query, args...); err != nil {
		return nil, err
	}

	return oh, nil
}

func (dbc *dbClient) BeginX() (*sqlx.Tx, error) {
	return dbc.db.Beginx()
}
