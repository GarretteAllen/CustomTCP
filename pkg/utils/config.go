package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	ServerPort    int    `mapstructure:"server_port"`
	DatabaseURI   string `mapstructure:"database_uri"`
	DatabaseName  string `mapstructure:"database_name"`
}

func LoadConfig(configPath, configName string) (*Config, error) {
	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %v", err)
	}

	return &config, nil
}
