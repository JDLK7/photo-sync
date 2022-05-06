package main

import (
	"photo-sync/internal/config"
	"photo-sync/internal/sync"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	logrus.WithField("config", cfg).Info("Loaded config")

	storage := sync.NewSQLiteStorage()
	syncService := sync.NewSyncService(cfg, storage)
	syncService.Listen()
}
