package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	DSN string
}

func Load() (*Config, error) {
	// Load .env file (KEY=VALUE) using viper without external libs
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("warning: %v", err)
	}

	// Get port from PORT (Render) or SERVER_PORT (local)
	port := viper.GetString("PORT")
	if port == "" {
		port = viper.GetString("SERVER_PORT")
	}
	if port == "" {
		port = "6000" // default port
	}

	// Get host from SERVER_HOST or default to 0.0.0.0
	host := viper.GetString("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	config := &Config{
		Server: ServerConfig{
			Port: port,
			Host: host,
		},
		Database: DatabaseConfig{
			DSN: viper.GetString("DATABASE_URL"),
		},
	}

	if config.Database.DSN == "" {
		return nil, fmt.Errorf("DATABASE_URL must be set")
	}

	return config, nil

}
