package postgres

import (
	"database/sql"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/table_api/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
)

type DBClient interface {
	DB() *sql.DB
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedQuery(query string, args interface{}) (*sqlx.Rows, error)
	GetTables(query string, args ...interface{}) ([]models.Table, error)
}

type dbClient struct {
	db *sqlx.DB
}

func NewDBClient(url *url.URL) (*dbClient, func() error, error) {
	driver := "postgres"

	db, err := sqlx.Open(driver, url.String())
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

func (dbc *dbClient) NamedQuery(query string, args interface{}) (*sqlx.Rows, error) {
	return dbc.db.NamedQuery(query, args)
}

func (dbc *dbClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dbc.db.Exec(query, args...)
}

func (dbc *dbClient) GetTables(query string, args ...interface{}) ([]models.Table, error) {
	tables := []models.Table{}
	if err := dbc.db.Select(&tables, query, args...); err != nil {
		return nil, err
	}
	return tables, nil
}
