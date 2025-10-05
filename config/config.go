package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		ChatReceiversCount int
		MessagesBufferSize int
		Address            string
	}
}

func New() *Config {

	return &Config{}
}

func (c *Config) LoadConfig() error {

	viper.SetConfigFile("temp/config.toml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("failed to unpack config: %v", err)
	}

	return nil
}
