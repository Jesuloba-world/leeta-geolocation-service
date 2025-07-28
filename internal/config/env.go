package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jesuloba-world/leeta-task/pkg/validator"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `json:"server" validate:"required"`
	Database DatabaseConfig `json:"database"`
	Storage  string         `json:"storage" validate:"required,oneof=memory postgres"`
}

type ServerConfig struct {
	Port         int `json:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  int `json:"read_timeout" validate:"required,min=1"`
	WriteTimeout int `json:"write_timeout" validate:"required,min=1"`
	IdleTimeout  int `json:"idle_timeout" validate:"required,min=1"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func LoadConfig() Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	config := Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 120),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "geolocation"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Storage: getEnv("STORAGE_TYPE", "memory"),
	}

	if err := ValidateConfig(config); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	return config
}

func ValidateConfig(cfg Config) error {
	if err := validator.ValidateStruct(cfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	if cfg.Storage == "postgres" {
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host is required when using postgres storage")
		}
		if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
			return fmt.Errorf("invalid database port: %d (must be between 1 and 65535)", cfg.Database.Port)
		}
		if cfg.Database.User == "" {
			return fmt.Errorf("database user is required when using postgres storage")
		}
		if cfg.Database.DBName == "" {
			return fmt.Errorf("database name is required when using postgres storage")
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
