package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string

	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDB       string
	PGSSLMode  string
}

func Load(path string) (Config, error) {
	_ = godotenv.Load(path)

	cfg := Config{
		HTTPAddr:   getenv("HTTP_ADDR", ":8080"),
		PGHost:     getenv("POSTGRES_HOST", "localhost"),
		PGPort:     getenv("POSTGRES_PORT", "5432"),
		PGUser:     getenv("POSTGRES_USER", "wallet"),
		PGPassword: getenv("POSTGRES_PASSWORD", "wallet"),
		PGDB:       getenv("POSTGRES_DB", "wallet"),
		PGSSLMode:  getenv("POSTGRES_SSLMODE", "disable"),
	}

	return cfg, nil
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PGUser, c.PGPassword, c.PGHost, c.PGPort, c.PGDB, c.PGSSLMode,
	)
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
