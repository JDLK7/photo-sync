package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	exclusions := []string{
		"/mnt/external-drive",
	}
	drives := []string{
		"/mnt/external-drive",
		"/mnt/sd-card",
		"/mnt/micro-sd-card",
	}

	expectedDrives := []string{
		"/mnt/sd-card",
		"/mnt/micro-sd-card",
	}
	assert.EqualValues(t, expectedDrives, Filter(exclusions, drives))
}
