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
		Name: "test_folder",
		Path: tempDir,
	}

	zipper := NewZIP()
	compressedFile, err := zipper.Compress(folder)
	assert.NoError(t, err)
	assert.FileExists(t, compressedFile.Id)

	zipReader, err := zip.OpenReader(compressedFile.Id)
	assert.NoError(t, err)
	defer zipReader.Close()

	os.Remove(compressedFile.Id)
}
