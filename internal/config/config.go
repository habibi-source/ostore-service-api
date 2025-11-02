// internal/config/config.go
package config

import "mini-project-ostore/pkg/database"

type Config struct {
	Server   ServerConfig
	Database database.Config
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		Database: database.Config{
			Host:     "localhost",
			Port:     "3306",
			User:     "ahmad",
			Password: "Habibi313367",
			Database: "ostore_db",
		},
		JWT: JWTConfig{
			Secret: "your-secret-key",
		},
	}
}
