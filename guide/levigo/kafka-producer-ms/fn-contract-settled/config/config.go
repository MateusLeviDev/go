package config

import (
	"flag"
	"fmt"
	"fn-contract-settled/pkg/constants"

	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Service config path")
}

type Config struct {
	Producer Producer `mapstructure:"producer"`
}

type Producer struct {
	BootstrapServers string
	Topic            string
	ClientID         string
}

func InitConfig() (*Config, error) {
	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return cfg, nil
}
