package fakeZIP

import (
	"archive/zip"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FakeZIP struct {
}

// Compress implements adapters.IVideoProcessorCompressor.
func (f *FakeZIP) Compress(file entities.Folder) (entities.File, error) {

	// Create the output zip file
	zipFile, err := os.Create(file.Path + ".zip")
	if err != nil {
		return entities.NewFile("", "", ""), err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the folder
	filepath.Walk(file.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path to maintain folder structure
		relativePath := strings.TrimPrefix(path, filepath.Dir(file.Path)+string(os.PathSeparator))
		if info.IsDir() {
			// Skip adding directories as zip adds them automatically
			return nil
		}

		// Create a new zip entry
		zipEntry, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		// Open the file to copy its content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy file content to the zip entry
		_, err = io.Copy(zipEntry, file)
		return err
	})

	log.Println("Mensagem publicada na fakeZIP - ", file)
	return entities.NewFile(filepath.Base(file.Path), zipFile.Name(), file.Path+".zip"), nil
}

func NewFakeZIP() adapters.IVideoProcessorCompressor {
	return &FakeZIP{}
}
