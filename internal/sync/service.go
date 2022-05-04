package sync

import (
	"os"
	"path/filepath"
	"photo-sync/internal/config"
	"photo-sync/pkg/utils/slices"
	"strings"
	"sync"
	"time"

	devices "github.com/deepakjois/gousbdrivedetector"
	"github.com/sirupsen/logrus"
)

var syncServiceInstance *SyncService
var syncServiceOnce sync.Once

func NewSyncService(cfg *config.Config, storage Storage) *SyncService {
	syncServiceOnce.Do(func() {
		syncServiceInstance = &SyncService{
			cfg:     cfg,
			storage: storage,
		}

		syncServiceInstance.init()
	})

	return syncServiceInstance
}

type SyncService struct {
	cfg            *config.Config
	storage        Storage
	currentDevices []string
}

func (s *SyncService) init() {
	var err error
	s.currentDevices, err = devices.Detect()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to detect new devices")
	}

	logrus.WithField("devices", s.currentDevices).Info("Devices detected at start")

	s.currentDevices = slices.Filter(s.cfg.ExcludedVolumes, s.currentDevices)
}

func (s *SyncService) findFiles(dir string) (files []string, err error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, de := range dirEntries {
		var subdirFiles []string
		path := filepath.Join(dir, de.Name())

		if de.IsDir() {
			subdirFiles, err = s.findFiles(path)
			if err != nil {
				return nil, err
			}

			files = append(files, subdirFiles...)
		} else {
			ext := strings.ToLower(filepath.Ext(path))
			if slices.Contains(s.cfg.SupportedExtensions, ext) {
				files = append(files, path)
			}
		}
	}

	return
}

func (s *SyncService) ScanDevice(device string) {
	logrus.WithField("device", device).Info("Scanning device")
	defer logrus.WithField("device", device).Info("Finished scanning device")

	files, err := s.findFiles(device)
	if err != nil {
		logrus.WithError(err).WithField("device", device).Error("Failed to find new files")
		return
	}

	logrus.WithField("device", device).WithField("files", files).Debug("Files found")
}

func (s *SyncService) Listen() {
	for {
		devices, err := devices.Detect()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to detect new devices")
		}

		devices = slices.Filter(s.cfg.ExcludedVolumes, devices)

		if len(devices) > len(s.currentDevices) {
			logrus.WithField("devices", devices).Info("New devices detected!")

			for _, device := range devices {
				s.ScanDevice(device)
			}
		}

		s.currentDevices = devices

		time.Sleep(time.Second * 10)
	}
}
