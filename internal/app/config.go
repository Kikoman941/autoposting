package app

import (
	ewrap "autoposting/pkg/err-wrapper"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PostgresDSN string
	IsProd      bool
	LogLevel    string
	ServerAddr  string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, ewrap.Errorf("cannot load env file: %w", err)
	}
	return &Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		IsProd:      os.Getenv("IS_PROD") == "true",
		LogLevel:    os.Getenv("LOG_LEVEL"),
		ServerAddr:  os.Getenv("SERVER_ADDR"),
	}, nil
}
