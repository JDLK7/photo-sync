package sync

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"photo-sync/internal/config"
	"photo-sync/internal/notification"
	"photo-sync/pkg/utils/slices"
	"strings"
	"sync"
	"time"

	devices "github.com/JDLK7/gousbdrivedetector"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const DateFormat = "2006-01-02"

var syncServiceInstance *SyncService
var syncServiceOnce sync.Once

func NewSyncService(cfg *config.Config, storage Storage) *SyncService {
	syncServiceOnce.Do(func() {
		syncServiceInstance = &SyncService{
			cfg:     cfg,
			storage: storage,
		}
	})

	return syncServiceInstance
}

// SyncService is responsible for listening to connected devices and
// synchronizing their new files onto another volume.
type SyncService struct {
	cfg            *config.Config
	storage        Storage
	currentDevices []string
}

// findFiles makes a recursive search for files in the given directory.
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

				// TODO: should this go here?
				s.ProcessFile(path)
			}
		}
	}

	return
}

// fileDigest calculates the sha256 digest of a file. This comes in handy for
// quickly comparing file contents.
func (s *SyncService) fileDigest(file io.Reader) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		logrus.WithError(err).WithField("file", file).Error("Failed to calculate file digest")
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// FIXME: for some reason the file might be left open
// processFile performs the heavy lifting of the synchronization:
// 1. Gets file's metadata
// 2. Creates a destination folder with the date when the file was created
// 3. Looks up the file to check if it was already process
// 4. Copies the file to the destination folder
// 5. Stores the file metadata to avoid reprocessing in the future
func (s *SyncService) processFile(srcPath string) error {
	srcStat, err := os.Stat(srcPath)
	if err != nil {
		return errors.Wrap(err, "Failed to stat source file")
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open source file")
	}
	defer srcFile.Close()

	folderName := srcStat.ModTime().Format(DateFormat)
	dstFolder := filepath.Join(s.cfg.DestinationVolume, folderName)

	err = os.MkdirAll(dstFolder, 0755)
	if err != nil {
		return errors.Wrapf(err, "Failed to create destination folder '%s'", dstFolder)
	}

	dstPath := filepath.Join(dstFolder, filepath.Base(srcPath))
	file, err := s.storage.Find(dstPath)
	if file != nil || err == nil {
		logrus.WithField("destination_file", dstPath).Warn("File already processed; skipping")
		return nil
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to create destination file '%s'", dstPath)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrapf(err, "Failed to copy file '%s'", dstPath)
	}

	err = s.storage.Create(&File{
		Name: dstPath,
		// FIXME: thos doesn't work and multiwriter kills the process with OOM
		// Digest: s.fileDigest(srcFile),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to mark file as processed")
	}

	return nil
}

// Process file just calls the unexported method 'processFile' and logs and
// notify any errors.
func (s *SyncService) ProcessFile(srcPath string) {
	err := s.processFile(srcPath)
	if err != nil {
		logrus.WithError(err).WithField("source_file", srcPath).Error("Failed to process file")
		notification.SendError(err)
	}
}

// logAndNotify logs and notifies a message related to the device.
func (s *SyncService) logAndNotify(device, message string) {
	logrus.WithField("device", device).Info(message)
	notification.SendMessagef("%s %s", message, device)
}

func (s *SyncService) ScanDevice(device string) {
	s.logAndNotify(device, "Scanning device")
	defer s.logAndNotify(device, "Finished scanning device")

	files, err := s.findFiles(device)
	if err != nil {
		logrus.WithError(err).WithField("device", device).Error("Failed to find new files")
		return
	}

	logrus.WithField("device", device).WithField("files", files).Debug("Files found")
}

// Listen is a blocking method that looks for new connected devices and
// automatically starts scanning them.
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
