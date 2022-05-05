package notification

import (
	"fmt"
	"photo-sync/internal/config"

	"github.com/gregdel/pushover"
	"github.com/sirupsen/logrus"
)

var app *pushover.Pushover
var recipient *pushover.Recipient
var enabled bool

// init loads the configuration and configures the pushover client if
// notifications are enabled.
func init() {
	cfg := config.NewConfig()

	enabled = cfg.Notifications.Enabled

	if cfg.Notifications.Enabled {
		app = pushover.New(cfg.Notifications.ApiKey)
		recipient = pushover.NewRecipient(cfg.Notifications.UserKey)
	}

}

// SendMessagef formats a message with the given arguments and sends a
// notification.
func SendMessagef(message string, args ...any) {
	SendMessage(fmt.Sprintf(message, args...))
}

// SendMessage sends a notification with the given message.
func SendMessage(message string) {
	if enabled {
		msg := pushover.NewMessageWithTitle(message, "Media Sync")
		_, err := app.SendMessage(msg, recipient)
		if err != nil {
			logrus.WithError(err).Error("Failed to deliver notification")
		}
	}
}

// SendError sends a notification formatted with the given error.
func SendError(err error) {
	if enabled {
		msg := pushover.NewMessageWithTitle(err.Error(), "Media Sync")
		_, err := app.SendMessage(msg, recipient)
		if err != nil {
			logrus.WithError(err).Error("Failed to deliver notification")
		}
	}
}
