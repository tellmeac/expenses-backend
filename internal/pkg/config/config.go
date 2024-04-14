package config

import (
	"github.com/spf13/viper"
	"strings"
)

// LoaderOption represents option to configure viper.Viper.
type LoaderOption func(*viper.Viper)

// WithConfigPath is an option to set configuration file path to lookup.
func WithConfigPath(p string) LoaderOption {
	return func(v *viper.Viper) {
		v.SetConfigFile(p)
	}
}

// PrepareLoader returns configured viper.Viper engine to use as config loader.
func PrepareLoader(options ...LoaderOption) *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	// Default env key replacing, can be overridden with options.
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	for _, opt := range options {
		opt(v)
	}

	return v
}
