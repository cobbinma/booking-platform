package config

import (
	"fmt"
	"os"
)

func Port() string {
	if p := os.Getenv("PORT"); p != "" {
		return fmt.Sprintf(":%s", p)
	}
	return ":8888"
}
