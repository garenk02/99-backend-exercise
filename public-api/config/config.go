package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort        string
	UserServiceURL    string
	ListingServiceURL string
}

// New returns a new Config with values from environment variables
func New() *Config {
	return &Config{
		ServerPort:        getEnvOrDefault("SERVER_PORT", "6002"),
		UserServiceURL:    getEnvOrDefault("USER_SERVICE_URL", "http://localhost:6001"),
		ListingServiceURL: getEnvOrDefault("LISTING_SERVICE_URL", "http://localhost:6000"),
	}
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
