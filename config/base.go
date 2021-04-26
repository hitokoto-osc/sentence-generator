package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// BuildTag is git commit tag
	BuildTag = "unknown"
	// BuildTime is program build time
	BuildTime = "unknown"
	// CommitTime is git commit Time
	CommitTime = "unknown"
	// Debug Mode
	Debug = false
	// Version is program version
	Version = "development"
)

// SetDebug change debug option
func SetDebug(v bool) {
	Debug = v
}

// InitConfigDriver init config
func InitConfigDriver(cfgFile string, logger *zap.Logger) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("data")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal(
			errors.WithMessage(
				err,
				"can't initialize config driver, program exit",
			).
				Error(),
		)
	} else {
		logger.Info("Using config file:" + viper.ConfigFileUsed())
	}
	if err := viper.UnmarshalKey("database", &Database); err != nil {
		logger.Fatal(
			errors.WithMessage(
				err,
				"can't reflect config data, program exit",
			).
				Error(),
		)
	}
	if err := viper.UnmarshalKey("core", &Core); err != nil {
		logger.Fatal(
			errors.WithMessage(
				err,
				"can't reflect config data, program exit",
			).
				Error(),
		)
	}
	if err := viper.UnmarshalKey("git", &Git); err != nil {
		logger.Fatal(
			errors.WithMessage(
				err,
				"can't reflect config data, program exit",
			).
				Error(),
		)
	}
}
