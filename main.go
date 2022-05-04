package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	devices "github.com/deepakjois/gousbdrivedetector"
)

var devicesAtStart []string

func listen() {
	for {
		mountPoints, err := devices.Detect()
		if err != nil {
			log.Fatalf("failed to detect new devices: %s", err.Error())
		}

		_ = mountPoints
	}
}

func findFiles(dir string) (files []string, err error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, de := range dirEntries {
		var subdirFiles []string
		path := filepath.Join(dir, de.Name())

		if de.IsDir() {
			subdirFiles, err = findFiles(path)
			if err != nil {
				return nil, err
			}

			files = append(files, subdirFiles...)
		} else {
			files = append(files, path)
		}
	}

	return
}

func init() {
	var err error
	devicesAtStart, err = devices.Detect()
	if err != nil {
		log.Fatalf("failed to detect new devices: %s", err.Error())
	}
}

func main() {
	files, err := findFiles(".")
	if err != nil {
		log.Fatalf("failed to find files: %s", err.Error())
	}

	fmt.Println(files)
}
