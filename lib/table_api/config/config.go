package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func Port() string {
	if p := os.Getenv("PORT"); p != "" {
		return fmt.Sprintf(":%s", p)
	}
	return ":8989"
}

func VenueAPIRoot() string {
	url := os.Getenv("VENUE_API_ROOT")
	if url == "" {
		logrus.Fatal("Venue API URL not set")
	}

	return url
}
