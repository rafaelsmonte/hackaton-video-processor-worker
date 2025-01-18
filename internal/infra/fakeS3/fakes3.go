package fakeS3

import (
	"fmt"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"io"
	"os"
	"path/filepath"
)

type FakeS3 struct {
}

func (f *FakeS3) Download(fi entities.File) (entities.File, error) {

	sourceFile, err := os.Open(fi.Path)
	if err != nil {
		return entities.File{}, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(fi.Name)
	if err != nil {
		return entities.File{}, fmt.Errorf("erro ao criar o arquivo de destino: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return entities.File{}, fmt.Errorf("erro ao copiar o arquivo: %v", err)
	}

	path, _ := os.Getwd()

	return entities.File{Path: path, Name: destinationFile.Name(), Id: destinationFile.Name()}, nil
}

func (f *FakeS3) Upload(fi entities.File) (string, error) {
	destinationDir := "/home/ra056172/fiap/hackaton-video-processor-worker/upload/"

	destinationPath := filepath.Join(destinationDir, fi.Id)

	sourceFile, err := os.Open(fi.Path)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath + ".zip")
	if err != nil {
		return "", fmt.Errorf("erro ao criar o arquivo de destino: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("erro ao copiar o arquivo: %v", err)
	}

	return "", nil
}

func NewS3() (adapters.IVideoProcessorStorage, error) {
	return &FakeS3{}, nil
}
