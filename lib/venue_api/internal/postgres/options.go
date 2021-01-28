package postgres

import (
	"github.com/google/uuid"
	"net/url"
)

func WithDatabaseURL(pgURL *url.URL) func(*client) {
	return func(c *client) {
		c.pgURL = pgURL
	}
}

func WithMigrationsSourceURL(url string) func(*client) {
	return func(c *client) {
		c.migrationsSource = url
	}
}

func WithStaticUUIDGenerator(id string) func(*client) {
	return func(c *client) {
		c.uuid = newStaticUUID(id)
	}
}

type uuidGenerator interface {
	UUID() string
}

var _ uuidGenerator = (*randomUUID)(nil)

type randomUUID struct{}

func newRandomUUID() *randomUUID {
	return &randomUUID{}
}

func (r randomUUID) UUID() string {
	return uuid.New().String()
}

var _ uuidGenerator = (*staticUUID)(nil)

type staticUUID struct {
	id string
}

func newStaticUUID(id string) *staticUUID {
	return &staticUUID{id: id}
}

func (s staticUUID) UUID() string {
	return s.id
}
