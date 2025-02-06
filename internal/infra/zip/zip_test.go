package zip

import (
	"archive/zip"
	"hackaton-video-processor-worker/internal/domain/entities"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_zip")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, "test.txt")
	err = ioutil.WriteFile(tempFilePath, []byte("Hello, World!"), 0644)
	assert.NoError(t, err)

	folder := entities.Folder{
		Name:   "test_folder",
		Path:   tempDir,
		Id:     "123",
		UserId: "12434",
	}

	zipper := NewZIP()
	compressedFile, err := zipper.Compress(folder)
	assert.NoError(t, err)
	assert.FileExists(t, compressedFile.Name)

	zipReader, err := zip.OpenReader(compressedFile.Name)
	assert.NoError(t, err)
	defer zipReader.Close()

	defer os.Remove(compressedFile.Name)
}
