package sync

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name   string
	Digest string
}
