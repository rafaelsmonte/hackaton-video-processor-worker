package zip

import (
	"archive/zip"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ZIP struct {
}

func (f *ZIP) Compress(file entities.Folder) (entities.File, error) {

	//tests de carga
	currentTime := time.Now().UnixMilli()
	currentTimeStr := strconv.FormatInt(currentTime, 10)
	//
	fileName := file.Name + "-" + currentTimeStr + ".zip"
	zipFile, err := os.Create(fileName)
	if err != nil {
		return entities.NewFile("", "", "", "", nil), err
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

	return entities.NewFile(file.Id, file.Path, file.UserId, fileName, nil), nil
}

func NewZIP() adapters.IVideoProcessorCompressor {
	return &ZIP{}
}
