package fakeStorage

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"log"
)

type FakeStorage struct {
}

func (f *FakeStorage) Download(file entities.File) (entities.File, error) {
	log.Println("Baixando o arquivo", file)
	return entities.NewFile(file.Id, file.Name, file.Path), nil
}

func (f *FakeStorage) Upload(file entities.File) error {
	log.Println("UPANDO o arquivo", file)
	return nil
}

func NewFakeStorage() adapters.IVideoProcessorStorage {
	return &FakeStorage{}
}
