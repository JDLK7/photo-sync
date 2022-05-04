package main

import (
	"os"
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

	actualFiles, err := findFiles("tmp")

	assert.NoError(t, err)
	assert.Equal(t, expectedFiles, actualFiles)
}
