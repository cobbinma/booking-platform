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
	return ":6969"
}

func TableAPIRoot() string {
	url := os.Getenv("TABLE_API_ROOT")
	if url == "" {
		logrus.Fatal("Table API URL not set")
	}

	return url
}

func GetAllowOrigin() string {
	return os.Getenv("ALLOW_ORIGIN")
}
