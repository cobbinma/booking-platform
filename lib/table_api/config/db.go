package config

import (
	"net/url"
	"os"
)

var (
	dbHost     = os.Getenv("DB_HOST")
	dbName     = os.Getenv("DB_NAME")
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbSSLMode  = os.Getenv("DB_SSLMODE")
)

func PostgresURL() *url.URL {
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Path:   dbName,
		Host:   dbHost,
	}
	q := pgURL.Query()
	q.Add("sslmode", dbSSLMode)
	pgURL.RawQuery = q.Encode()
	return pgURL
}
