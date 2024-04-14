package config

import "github.com/spf13/viper"

type Config struct {
	ListenPort         string
	DatabaseConnection string
}

func ParseConfig(loader *viper.Viper) (*Config, error) {
	cfg := &Config{}

	if err := loader.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := loader.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
