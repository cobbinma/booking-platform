package postgres

import "net/url"

func WithDatabaseURL(pgURL *url.URL) func(*client) {
	return func(c *client) {
		c.pgURL = pgURL
	}
}
