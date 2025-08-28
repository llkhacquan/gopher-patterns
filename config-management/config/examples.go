package config

import (
	"github.com/pkg/errors"
)

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Addresses []string `mapstructure:"addresses"`
}

// TradingConfig holds trading-specific settings
type TradingConfig struct {
	MaxOrdersPerUser int `mapstructure:"max_orders_per_user"`
}

// AppConfig represents the main application configuration
type AppConfig struct {
	ServiceName string         `mapstructure:"service_name"`
	Database    DatabaseConfig `mapstructure:"database"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Trading     TradingConfig  `mapstructure:"trading"`
}

// Init initializes configuration using the simple pattern
func Init() (AppConfig, error) {
	InitViper()
	var cfg AppConfig
	if err := Unmarshal(&cfg); err != nil {
		return AppConfig{}, errors.Wrap(err, "failed to unmarshal config")
	}
	return cfg, nil
}

// MustInit initializes configuration and panics on error
func MustInit() AppConfig {
	cfg, err := Init()
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize config"))
	}
	return cfg
}
