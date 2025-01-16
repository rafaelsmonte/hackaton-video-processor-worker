package usecases

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"path/filepath"
)

type IConvertVideoUsecase interface {
	Execute(ConvertVideoInput ConvertVideoInput) (ConvertVideoOutput, error)
}

type ConvertVideoUsecase struct {
	videoProcessorConverter  adapters.IVideoProcessorConverter
	videoProcessorMessaging  adapters.IVideoProcessorMessaging
	videoProcessorStorage    adapters.IVideoProcessorStorage
	videoProcessorCompressor adapters.IVideoProcessorCompressor
}

type ConvertVideoInput struct {
	VideoName string
	VideoPath string
	VideoId   string
}

type ConvertVideoOutput struct {
	VideoPath string
}

func NewConvertVideoUsecase(
	videoProcessorConverter adapters.IVideoProcessorConverter,
	videoProcessorMessaging adapters.IVideoProcessorMessaging,
	videoProcessorStorage adapters.IVideoProcessorStorage,
	videoProcessorCompressor adapters.IVideoProcessorCompressor,
) ConvertVideoUsecase {
	return ConvertVideoUsecase{
		videoProcessorConverter:  videoProcessorConverter,
		videoProcessorMessaging:  videoProcessorMessaging,
		videoProcessorStorage:    videoProcessorStorage,
		videoProcessorCompressor: videoProcessorCompressor,
	}
}

var workerPool = make(chan struct{}, 5)

func (u *ConvertVideoUsecase) Execute(ConvertVideoInput ConvertVideoInput) (ConvertVideoOutput, error) {

	go func() {
		workerPool <- struct{}{}
		defer func() { <-workerPool }()

		//TODO  Messages
		//		Deletar arquivos no disco
		err := u.videoProcessorMessaging.Publish("Processando")
		if err != nil {
		}

		sourceFile := entities.NewFile(ConvertVideoInput.VideoId, ConvertVideoInput.VideoName, ConvertVideoInput.VideoPath)

		newFile, err := u.videoProcessorStorage.Download(sourceFile)
		if err != nil {
			u.videoProcessorMessaging.Publish("Erro ao baixar video " + err.Error())
			return
		}
		defer os.Remove(filepath.Join(newFile.Path, newFile.Name))

		newFolder, err := u.videoProcessorConverter.ConvertToImages(newFile)
		defer os.RemoveAll(newFolder.Path)
		if err != nil {
			u.videoProcessorMessaging.Publish("Erro ao converter video " + err.Error())
			return
		}
		compressedFile, err := u.videoProcessorCompressor.Compress(newFolder)
		if err != nil {
			u.videoProcessorMessaging.Publish("Erro ao zipar video " + err.Error())
			return
		}
		err = u.videoProcessorStorage.Upload(compressedFile)
		if err != nil {
			u.videoProcessorMessaging.Publish("Erro ao upar video " + err.Error())
			return
		}
		defer os.Remove(compressedFile.Path + "")
		u.videoProcessorMessaging.Publish("Complete")

	}()

	return ConvertVideoOutput{}, nil
}
