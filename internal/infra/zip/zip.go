package zip

import (
	"archive/zip"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZIP struct {
}

func (f *ZIP) Compress(file entities.Folder) (entities.File, error) {
	zipFile, err := os.Create(file.Path + ".zip")
	if err != nil {
		return entities.NewFile("", "", ""), err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	filepath.Walk(file.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path, filepath.Dir(file.Path)+string(os.PathSeparator))
		if info.IsDir() {
			return nil
		}

		zipEntry, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})
	return entities.NewFile(filepath.Base(file.Path), zipFile.Name(), file.Path+".zip"), nil
}

func NewZIP() adapters.IVideoProcessorCompressor {
	return &ZIP{}
}
