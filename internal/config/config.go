package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PostgresqlDSN string
	IsProd        bool
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Cannot load env file: %v", err)
	}
	return &Config{
		PostgresqlDSN: os.Getenv("POSTGRESQL_DSN"),
		IsProd:        os.Getenv("IS_PROD") == "true",
	}
}
