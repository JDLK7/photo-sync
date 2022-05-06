package sync

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDigest(t *testing.T) {
	content1 := []byte("Nulla irure eu amet pariatur duis.")
	content2 := []byte("Commodo nostrud est id excepteur ea.")

	file1, _ := os.Create("file1")
	file1.Write(content1)
	defer os.Remove("file1")
	defer file1.Close()

	file2, _ := os.Create("file2")
	file2.Write(content2)
	defer os.Remove("file2")
	defer file2.Close()

	svc := &SyncService{}

	file1, _ = os.Open("file1")
	file2, _ = os.Open("file2")

	digest1 := svc.fileDigest(file1)
	digest2 := svc.fileDigest(file2)

	assert.NotEqual(t, digest1, digest2)
}
