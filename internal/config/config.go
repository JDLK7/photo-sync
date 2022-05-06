package config

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NotificationsConfig holds the configuration for Pushover notifications.
type NotificationsConfig struct {
	// Enabled controls whether notifications should be sent or not.
	Enabled bool
	// UserKey is the key of the notification recipient.
	UserKey string
	// ApiKey is the API key of the application that will send the
	// notifications.
	ApiKey string
}

// Config holds the configuration of the application.
type Config struct {
	// DestinationVolume is the path to the volume where the synchronized files
	// will be stored.
	DestinationVolume string
	// ExcludedVolumes is a list of paths to volumes that should not be synched.
	ExcludedVolumes []string
	// SupportedExtensions is a list of extensions that should be synched. They
	// must start with a dot.
	SupportedExtensions []string
	// Notifications holds the configuration for Pushover notifications.
	Notifications NotificationsConfig
}

var config *Config
var once sync.Once

// NewConfig initializes the config singleton by loading the configuration file
// and unmarshalling it.
func NewConfig() *Config {
	once.Do(func() {
		config = &Config{}

		viper.AddConfigPath(".")
		if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
			viper.AddConfigPath(configPath)
		}

		err := viper.ReadInConfig()
		if err != nil {
			logrus.WithError(err).Error("Failed to read configuration")
			return
		}

		err = viper.Unmarshal(config)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to unmarshal configuration")
		}
	})

	return config
}
