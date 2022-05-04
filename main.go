package main

import (
	"os"
	"path/filepath"
	"photo-sync/internal/config"
	"photo-sync/pkg/utils/slices"
	"strings"
	"time"

	devices "github.com/deepakjois/gousbdrivedetector"
	"github.com/sirupsen/logrus"
)

var currentDevices []string

func listen(cfg *config.Config) {
	for {
		devices, err := devices.Detect()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to detect new devices")
		}

		devices = slices.Filter(cfg.ExcludedVolumes, devices)

		if len(devices) > len(currentDevices) {
			logrus.WithField("devices", devices).Info("New devices detected!")

			for _, device := range devices {
				logrus.WithField("device", device).Info("Scanning device")

				files, err := findFiles(cfg, device)
				if err != nil {
					logrus.WithError(err).WithField("device", device).Error("Failed to find new files")
					break
				}

				logrus.WithField("device", device).WithField("files", files).Debug("Files found")
			}
		}

		currentDevices = devices

		time.Sleep(time.Second * 10)
	}
}

func findFiles(cfg *config.Config, dir string) (files []string, err error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, de := range dirEntries {
		var subdirFiles []string
		path := filepath.Join(dir, de.Name())

		if de.IsDir() {
			subdirFiles, err = findFiles(cfg, path)
			if err != nil {
				return nil, err
			}

			files = append(files, subdirFiles...)
		} else {
			ext := strings.ToLower(filepath.Ext(path))
			if slices.Contains(cfg.SupportedExtensions, ext) {
				files = append(files, path)
			}
		}
	}

	return
}

func main() {
	cfg := config.NewConfig()
	logrus.WithField("config", cfg).Info("Loaded config")

	var err error
	currentDevices, err = devices.Detect()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to detect new devices")
	}

	logrus.WithField("devices", currentDevices).Info("Devices detected at start")

	currentDevices = slices.Filter(cfg.ExcludedVolumes, currentDevices)

	listen(cfg)
}
