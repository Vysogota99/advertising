package server

import (
	"fmt"
	"os"
)

// Config ...
type Config struct {
	serverPort   string
	dbConnString string
}

// NewConfig - helper to init config
func NewConfig() (*Config, error) {
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		return nil, fmt.Errorf("No SERVER_PORT in .env")
	}

	dbConnString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		return nil, fmt.Errorf("No DB_CONN_STRING in .env")
	}

	return &Config{
		serverPort:   serverPort,
		dbConnString: dbConnString,
	}, nil
}
