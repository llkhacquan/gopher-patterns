# Config Management Pattern

A comprehensive configuration management pattern using Viper with support for YAML files, environment variable overrides, modular config loading, and custom unmarshaling.

## Key Features

- **YAML Configuration**: Primary config loaded from YAML files
- **Environment Overrides**: Environment variables override YAML values
- **Modular Configuration**: Additional configs pattern for service-specific settings
- **Small Config Structs**: Focused, maintainable configuration structures
- **Custom Unmarshaling**: Support for complex data types and validation
- **Thread-Safe**: Safe for concurrent access across goroutines
- **Validation**: Built-in validation with meaningful error messages

## Architecture

```
config/
├── config.go          # Core configuration loader and types
├── database.go        # Database configuration module
├── server.go          # HTTP server configuration module
├── logging.go         # Logging configuration module
└── additional.go      # Additional configs pattern implementation
```

## Configuration Structure

### Main Config
```yaml
# config.yaml
app:
  name: "spotx-exchange"
  environment: "development"
  debug: true

database:
  host: "localhost"
  port: 5432
  name: "spotx"
  user: "postgres"
  password: "password"
  max_connections: 100
  ssl_mode: "disable"

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  shutdown_timeout: "10s"

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

### Environment Overrides
```bash
# Environment variables override YAML values
export APP_NAME="spotx-production"
export DATABASE_HOST="prod-db.example.com"
export SERVER_PORT="8443"
export LOGGING_LEVEL="warn"
```

### Additional Configs (Modular)
```yaml
# configs/trading.yaml
trading:
  max_orders_per_user: 1000
  order_timeout: "5m"
  price_precision: 8

# configs/monitoring.yaml  
monitoring:
  metrics_enabled: true
  metrics_port: 9090
  health_check_interval: "30s"
```

## Usage Examples

### Basic Configuration Loading
```go
package main

import (
    "log"
    "github.com/xtrading/nova/bin/gopher-patterns/config-management/config"
)

func main() {
    // Load main configuration
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Use configuration
    log.Printf("Starting %s in %s mode", cfg.App.Name, cfg.App.Environment)
    
    // Database config
    db, err := setupDatabase(cfg.Database)
    if err != nil {
        log.Fatal("Database setup failed:", err)
    }
    
    // Server config
    server := setupServer(cfg.Server, db)
    server.Start()
}
```

### Modular Configuration Loading
```go
// Load additional configs
tradingCfg, err := config.LoadAdditional[TradingConfig]("configs/trading.yaml")
if err != nil {
    log.Fatal("Failed to load trading config:", err)
}

monitoringCfg, err := config.LoadAdditional[MonitoringConfig]("configs/monitoring.yaml")
if err != nil {
    log.Fatal("Failed to load monitoring config:", err)
}

// Use modular configs
tradingService := NewTradingService(tradingCfg)
monitoringService := NewMonitoringService(monitoringCfg)
```

### Custom Unmarshaling
```go
type DatabaseConfig struct {
    Host            string        `yaml:"host" validate:"required"`
    Port            int           `yaml:"port" validate:"min=1,max=65535"`
    Name            string        `yaml:"name" validate:"required"`
    User            string        `yaml:"user" validate:"required"`
    Password        string        `yaml:"password" validate:"required"`
    MaxConnections  int           `yaml:"max_connections" validate:"min=1"`
    SSLMode         string        `yaml:"ssl_mode" validate:"oneof=disable require verify-ca verify-full"`
    ConnectTimeout  time.Duration `yaml:"connect_timeout"`
}

// Custom unmarshaling for duration fields
func (d *DatabaseConfig) UnmarshalYAML(value *yaml.Node) error {
    // Custom parsing logic for complex fields
    // Validation logic
    // Default value assignment
}
```

### Small Config Structs Pattern
```go
// Instead of one massive config struct, use focused structs
type ServerConfig struct {
    Host            string        `yaml:"host"`
    Port            int           `yaml:"port"`
    ReadTimeout     time.Duration `yaml:"read_timeout"`
    WriteTimeout    time.Duration `yaml:"write_timeout"`
    ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type LoggingConfig struct {
    Level  string `yaml:"level" validate:"oneof=debug info warn error"`
    Format string `yaml:"format" validate:"oneof=json text"`
    Output string `yaml:"output"`
}

// Easy to pass around and test
func NewHTTPServer(cfg ServerConfig) *http.Server {
    return &http.Server{
        Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    }
}
```

## Environment Variable Mapping

The pattern automatically maps environment variables to YAML keys using the following rules:

- Nested keys: `database.host` → `DATABASE_HOST`
- Arrays: `database.replicas[0].host` → `DATABASE_REPLICAS_0_HOST`
- Case conversion: All uppercase for environment variables
- Delimiter: Underscore (`_`) separates nested levels

## Validation

Configuration validation happens automatically during loading:

```go
type Config struct {
    App      AppConfig      `yaml:"app" validate:"required"`
    Database DatabaseConfig `yaml:"database" validate:"required"`
    Server   ServerConfig   `yaml:"server" validate:"required"`
}

// Validation tags are enforced
// Custom validation methods can be implemented
// Meaningful error messages are provided
```

## Best Practices

1. **Small Structs**: Keep configuration structs focused and small
2. **Validation**: Always validate configuration values
3. **Defaults**: Provide sensible defaults for optional values
4. **Environment Overrides**: Use environment variables for deployment-specific values
5. **Modular Loading**: Use additional configs for service-specific settings
6. **Type Safety**: Use proper Go types (time.Duration, url.URL, etc.)
7. **Documentation**: Document all configuration options clearly

## Testing

```go
func TestConfigLoading(t *testing.T) {
    // Test YAML loading
    cfg, err := config.LoadFromString(yamlContent)
    assert.NoError(t, err)
    
    // Test environment overrides
    os.Setenv("DATABASE_HOST", "test-db")
    cfg, err = config.Load("test-config.yaml")
    assert.Equal(t, "test-db", cfg.Database.Host)
    
    // Test validation
    invalidCfg := Config{Database: DatabaseConfig{Port: 0}}
    err = config.Validate(invalidCfg)
    assert.Error(t, err)
}
```

## Two Implementation Approaches

This pattern provides two different implementation approaches:

### 1. Standard Approach (`example/`)
- Direct YAML unmarshaling with custom types
- Validator-based validation 
- Generic additional configs loading
- Best for: Simple to medium complexity applications
- **Note**: Custom unmarshaling only works for additional configs, not main config due to Viper limitations

### 2. Advanced Approach (`example-advanced/`)
- Viper-based configuration with decode hooks
- mapstructure tags for field mapping
- Environment-based config files (config.local.yaml, config.dev.yaml, etc.)
- additional_configs pattern for modular loading
- Custom decode hooks for decimal.Decimal, time.Duration, etc.
- Methods on config structs for initialization
- Best for: Complex microservices architectures

## Usage Examples

### Standard Approach
```bash
make example          # Run standard example
make test            # Run all tests
```

### Advanced Approach
```bash
make example-advanced # Run advanced example with decode hooks
```

## Files Generated

### Standard Approach
- `config/config.go` - Core configuration loader and main config struct
- `config/database.go` - Database configuration module with custom unmarshaling
- `config/server.go` - HTTP server configuration module
- `config/logging.go` - Logging configuration module
- `config/additional.go` - Additional configs pattern implementation
- `example/main.go` - Standard usage example
- `configs/` - Example YAML configuration files

### Advanced Approach
- `shared/config/config.go` - Viper initialization with additional_configs pattern
- `shared/config/unmarshal.go` - Custom decode hooks for complex types
- `shared/config/examples.go` - Small, focused config structs with methods
- `example-advanced/main.go` - Advanced usage example
- `config.local.yaml` - Environment-based config file
- `config-additional/` - Modular configuration files

This pattern provides a robust, maintainable configuration management system that scales from simple applications to complex microservices architectures.