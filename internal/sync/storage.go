package sync

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage interface {
	Create(file *File) error
	Find(name string) (*File, error)
}

func NewSQLiteStorage() *SQLStorage {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open SQLite connection")
	}

	storage := &SQLStorage{db: db}
	err = storage.Migrate()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to execute SQLite migrations")
	}

	return storage
}

type SQLStorage struct {
	db *gorm.DB
}

func (s SQLStorage) Migrate() error {
	err := s.db.AutoMigrate(&File{})
	return err
}

func (s SQLStorage) Create(file *File) error {
	return s.db.Create(file).Error
}

func (s SQLStorage) Find(name string) (*File, error) {
	file := &File{}

	err := s.db.First(file, File{Name: name}).Error
	if err != nil {
		return nil, err
	}

	return file, nil
}
