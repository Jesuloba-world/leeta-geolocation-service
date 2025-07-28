package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default configuration
	cfg := LoadConfig()

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Server.ReadTimeout != 10 {
		t.Errorf("Expected default read timeout 10, got %d", cfg.Server.ReadTimeout)
	}

	if cfg.Storage != "memory" {
		t.Errorf("Expected default storage 'memory', got %s", cfg.Storage)
	}
}

func TestLoadConfigWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("SERVER_READ_TIMEOUT", "15")
	os.Setenv("STORAGE_TYPE", "postgres")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")

	// Clean up after test
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_READ_TIMEOUT")
		os.Unsetenv("STORAGE_TYPE")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_NAME")
	}()

	cfg := LoadConfig()

	if cfg.Server.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", cfg.Server.Port)
	}

	if cfg.Server.ReadTimeout != 15 {
		t.Errorf("Expected read timeout 15, got %d", cfg.Server.ReadTimeout)
	}

	if cfg.Storage != "postgres" {
		t.Errorf("Expected storage 'postgres', got %s", cfg.Storage)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("Expected database host 'testhost', got %s", cfg.Database.Host)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid memory config",
			config: Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  10,
					WriteTimeout: 10,
					IdleTimeout:  120,
				},
				Storage: "memory",
			},
			wantErr: false,
		},
		{
			name: "valid postgres config",
			config: Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  10,
					WriteTimeout: 10,
					IdleTimeout:  120,
				},
				Database: DatabaseConfig{
					Host:   "localhost",
					Port:   5432,
					User:   "user",
					DBName: "db",
				},
				Storage: "postgres",
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: Config{
				Server: ServerConfig{
					Port:         0,
					ReadTimeout:  10,
					WriteTimeout: 10,
					IdleTimeout:  120,
				},
				Storage: "memory",
			},
			wantErr: true,
		},
		{
			name: "invalid storage type",
			config: Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  10,
					WriteTimeout: 10,
					IdleTimeout:  120,
				},
				Storage: "invalid",
			},
			wantErr: true,
		},
		{
			name: "postgres config missing host",
			config: Config{
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  10,
					WriteTimeout: 10,
					IdleTimeout:  120,
				},
				Database: DatabaseConfig{
					Host: "",
					Port: 5432,
				},
				Storage: "postgres",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got %s", result)
	}

	// Test with non-existing environment variable
	result = getEnv("NON_EXISTING_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got %s", result)
	}

	// Test with empty environment variable
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result = getEnv("EMPTY_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default' for empty var, got %s", result)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getEnvAsInt("TEST_INT", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid integer
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	result = getEnvAsInt("TEST_INVALID_INT", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}

	// Test with non-existing environment variable
	result = getEnvAsInt("NON_EXISTING_INT", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}
}