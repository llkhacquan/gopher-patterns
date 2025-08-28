package config

import (
	"testing"
)

func TestInitViper(t *testing.T) {
	// Set up test environment
	t.Setenv("RUNTIME_ENV", "local")

	// This should not panic
	InitViper()

	// Test that we can unmarshal to our config struct
	var cfg AppConfig
	err := Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify basic config values
	if cfg.ServiceName != "config_demo" {
		t.Errorf("Expected service_name 'config_demo', got %s", cfg.ServiceName)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected database host 'localhost', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected database port 5432, got %d", cfg.Database.Port)
	}

	if len(cfg.Redis.Addresses) == 0 || cfg.Redis.Addresses[0] != "localhost:6379" {
		t.Errorf("Expected redis address 'localhost:6379', got %v", cfg.Redis.Addresses)
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	t.Setenv("RUNTIME_ENV", "local")
	t.Setenv("SERVICE_NAME", "test-service")
	t.Setenv("DATABASE_HOST", "test-db")
	t.Setenv("DATABASE_PORT", "3306")

	InitViper()

	var cfg AppConfig
	err := Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify environment overrides work
	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected service_name 'test-service', got %s", cfg.ServiceName)
	}

	if cfg.Database.Host != "test-db" {
		t.Errorf("Expected database host 'test-db', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 3306 {
		t.Errorf("Expected database port 3306, got %d", cfg.Database.Port)
	}
}

func TestInit(t *testing.T) {
	t.Setenv("RUNTIME_ENV", "local")

	cfg, err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify config is not empty
	if cfg.ServiceName == "" {
		t.Error("Expected service_name to be set")
	}
	if cfg.Database.Host == "" {
		t.Error("Expected database host to be set")
	}
	if cfg.Database.Port == 0 {
		t.Error("Expected database port to be set")
	}
	if len(cfg.Redis.Addresses) == 0 {
		t.Error("Expected redis addresses to be set")
	}
	if cfg.Trading.MaxOrdersPerUser == 0 {
		t.Error("Expected trading config to be set")
	}
}

func TestMustInit(t *testing.T) {
	t.Setenv("RUNTIME_ENV", "local")

	// This should not panic with valid config
	cfg := MustInit()

	// Verify config is not empty
	if cfg.ServiceName == "" {
		t.Error("Expected service_name to be set")
	}
	if cfg.Database.Host == "" {
		t.Error("Expected database host to be set")
	}
	if cfg.Database.Port == 0 {
		t.Error("Expected database port to be set")
	}
	if len(cfg.Redis.Addresses) == 0 {
		t.Error("Expected redis addresses to be set")
	}
	if cfg.Trading.MaxOrdersPerUser == 0 {
		t.Error("Expected trading config to be set")
	}
}

func TestInitViperWithUnmarshal(t *testing.T) {
	t.Setenv("RUNTIME_ENV", "local")

	// Call InitViper then unmarshal
	InitViper()

	var cfg AppConfig
	err := Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Expect the config is loaded (not empty)
	if cfg.ServiceName == "" {
		t.Error("Config is empty: service_name not set")
	}
	if cfg.Database.Host == "" {
		t.Error("Config is empty: database host not set")
	}
	if cfg.Database.Port == 0 {
		t.Error("Config is empty: database port not set")
	}
	if len(cfg.Redis.Addresses) == 0 {
		t.Error("Config is empty: redis addresses not set")
	}
	if cfg.Trading.MaxOrdersPerUser == 0 {
		t.Error("Config is empty: trading config not set")
	}
}
