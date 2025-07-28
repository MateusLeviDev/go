package config

import (
	"flag"
	"fmt"

	"github.com/MateusLeviDev/pkg/constants"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Service config path")
}

type Config struct {
	Http      Http      `mapstructure:"http"`
	Redis     Redis     `mapstructure:"redis"`
	Processor Processor `mapstructure:"processor"`
	Workers   int       `mapstructure:"workers"`
}

type Http struct {
	Port         string `mapstructure:"port"`
	PaymentsPath string `mapstructure:"paymentsPath"`
	SummaryPath  string `mapstructure:"summaryPath"`
	HealthPath   string `mapstructure:"healthPath"`
}

type Redis struct {
	Address string `mapstructure:"address"`
}

type Processor struct {
	DefaultURL  string `mapstructure:"defaultUrl"`
	FallbackURL string `mapstructure:"fallbackUrl"`
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
