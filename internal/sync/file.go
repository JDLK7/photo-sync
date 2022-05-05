package sync

import "gorm.io/gorm"

// File is a simple database model that represents an already synched file.
type File struct {
	gorm.Model
	Name   string
	Digest string
}
