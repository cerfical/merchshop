package config

import (
	"errors"
	"strings"

	"github.com/cerfical/merchshop/internal/httpserv"
	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/postgres"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// MustLoad either successfully loads the configuration with [Load], or fails and causes the application to exit with an error.
func MustLoad(args []string) *Config {
	cfg, err := Load(args)
	if err != nil {
		log := log.New(&log.Config{})
		log.Fatal("Failed to load the configuration", err)
	}
	return cfg
}

// Load loads the configuration from environment variables or from the command line arguments specified.
func Load(args []string) (*Config, error) {
	v := viper.New()
	if len(args) > 1 {
		if len(args) != 2 {
			return nil, errors.New("expected a config path as the only command line argument")
		}
		v.SetConfigFile(args[1])
	}
	return load(v)
}

func load(v *viper.Viper) (*Config, error) {
	// Set up automatic configuration loading from environment variables of the same name
	// Build tag viper_bind_struct is required to properly unmarshal into a struct
	// TODO: https://github.com/spf13/viper/issues/1797
	v.SetEnvPrefix("merchshop")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// Make the configuration file optional
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Set up defaults
	v.SetDefault("log.level", log.LevelInfo)
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", "8080")

	// Set up defaults commonly used by Postgres
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", "5432")
	v.SetDefault("db.name", "postgres")
	v.SetDefault("db.user", "postgres")

	options := []viper.DecoderConfigOption{
		// Apply a custom hook so that [log.Level] values can be decoded with [log.Level.UnmarshalText]
		viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()),

		// Disallow unknown field names
		func(dc *mapstructure.DecoderConfig) {
			dc.ErrorUnused = true
		},
	}

	var cfg Config
	if err := v.Unmarshal(&cfg, options...); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Config encompasses all available application-level configuration settings.
type Config struct {
	Server httpserv.Config
	DB     postgres.Config
	Log    log.Config
}
