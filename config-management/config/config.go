package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	// runtime.Caller(0) gives us the current file path (/path/to/your-repo/config/config.go)
	// filepath.Dir(b) gives us the directory (/path/to/your-repo/config/)
	// We use ".." to go up one level to reach the project root (/path/to/your-repo/)
	// Adjust the number of "../" based on how deep your config package is nested
	Root = filepath.Join(filepath.Dir(b), "..")
)

// InitViper initializes Viper configuration with environment-based config loading
// It looks for config files named config.{RUNTIME_ENV}.yaml (e.g., config.local.yaml, config.prod.yaml)
// and supports additional config files through the additional_configs pattern
func InitViper(configPaths ...string) {
	// Determine environment (defaults to "local" if RUNTIME_ENV not set)
	env := os.Getenv("RUNTIME_ENV")
	if env == "" {
		env = "local"
	}

	// Look for config.{env}.yaml files
	viper.SetConfigName(fmt.Sprintf("config.%s", env))

	// Add custom config paths if provided
	for _, cp := range configPaths {
		// Join with Root so we can run app from any directory
		viper.AddConfigPath(path.Join(Root, cp))
	}

	// Add standard config search paths
	viper.AddConfigPath(".")                        // Current directory
	viper.AddConfigPath("./config")                 // ./config/ directory
	viper.AddConfigPath("./configs")                // ./configs/ directory
	viper.AddConfigPath(path.Join(Root, "configs")) // Project root configs/ directory

	// Load the main config file
	if err := viper.MergeInConfig(); err != nil {
		zap.L().Fatal("can't load config", zap.Error(err))
	}

	// Load additional config files specified in additional_configs array
	if err := loadAdditionalConfigs(Root); err != nil {
		zap.L().Fatal("can't load additional config", zap.Error(err))
	}

	// Enable automatic environment variable binding
	// This allows DATABASE_HOST env var to override database.host config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Merge environment variables with config
	if err := viper.MergeInConfig(); err != nil {
		zap.L().Fatal("can't merge config with env var", zap.Error(err))
	}
}

// loadAdditionalConfigs loads additional configuration files specified in the main config
// This pattern allows you to split configuration into multiple files for better organization
// Example: additional_configs: ["./shared.yaml", "./secrets.yaml"]
func loadAdditionalConfigs(configDir string) error {
	configFiles := viper.GetStringSlice("additional_configs")
	for _, file := range configFiles {
		abs, err := filepath.Abs(path.Join(configDir, file))
		if err != nil {
			return errors.Wrapf(err, "can't get absolute path for %s", file)
		}
		viper.SetConfigFile(abs)
		if err := viper.MergeInConfig(); err != nil {
			return errors.Wrapf(err, "can't load config file: %s", abs)
		}
	}
	return nil
}

// Unmarshal unmarshals the configuration into the provided struct
func Unmarshal(c any) error {
	if err := viper.Unmarshal(&c); err != nil {
		return errors.Wrap(err, "failed when unmarshal config")
	}
	return nil
}
