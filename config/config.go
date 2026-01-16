package config

import (
	"os"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Path string
}

type JWTConfig struct {
	Secret []byte
}

type ServerConfig struct {
	Port string
}

var AppConfig *Config

func LoadConfig() error {
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
	return nil
}

func (c *Config) GetDatabasePath() string {
	return c.Database.Path
}

func (c *Config) GetServerPort() string {
	if c.Server.Port == "" {
		return ":8080"
	}
	if c.Server.Port[0] != ':' {
		return ":" + c.Server.Port
	}
	return c.Server.Port
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
