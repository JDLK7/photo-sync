package main

import (
	"os"
	"photo-sync/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTests() {
	os.Create("tmp/metadata.txt")
	os.MkdirAll("tmp/SDCARD/DCIM/", 0755)
	os.Create("tmp/SDCARD/DCIM/20220428_0001.jpeg")
	os.MkdirAll("tmp/SDCARD/MEDIA/", 0755)
	os.Create("tmp/SDCARD/MEDIA/20220428_0001.mov")
}

func teardownTests() {
	os.RemoveAll("tmp")
}

func TestFindFiles(t *testing.T) {
	setupTests()
	defer teardownTests()

	expectedFiles := []string{
		"tmp/metadata.txt",
		"tmp/SDCARD/DCIM/20220428_0001.jpeg",
		"tmp/SDCARD/MEDIA/20220428_0001.mov",
	}
	cfg := &config.Config{
		SupportedExtensions: []string{"txt", "jpeg", "mov"},
	}

	actualFiles, err := findFiles(cfg, "tmp")

	assert.NoError(t, err)
	assert.Equal(t, expectedFiles, actualFiles)
}
