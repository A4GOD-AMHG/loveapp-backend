package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string // SQLite database file path
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret []byte
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

var AppConfig *Config

// InitConfig initializes the application configuration
func InitConfig() {
	AppConfig = &Config{
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./data/loveapp.db"),
		},
		JWT: JWTConfig{
			Secret: []byte(getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

// GetDatabasePath returns the database file path
func (c *Config) GetDatabasePath() string {
	return c.Database.Path
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
