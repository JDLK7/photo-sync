package main

import (
	"photo-sync/internal/config"
	"photo-sync/internal/sync"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	logrus.WithField("config", cfg).Info("Loaded config")

	syncService := sync.NewSyncService(cfg, nil)
	syncService.Listen()
}
