package config

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ExcludedVolumes     []string
	SupportedExtensions []string
}

var config *Config
var once sync.Once

func NewConfig() *Config {
	once.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
			viper.AddConfigPath(configPath)
		}

		err := viper.ReadInConfig()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to read configuration")
		}

		err = viper.Unmarshal(config)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to unmarshal configuration")
		}
	})

	return config
}