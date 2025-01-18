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
	VideoName        string
	VideoPath        string
	VideoId          string
	UserId           string
	VideoDescription string
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

func (converVideo *ConvertVideoUsecase) Execute(ConvertVideoInput ConvertVideoInput) (ConvertVideoOutput, error) {

	go func() {
		workerPool <- struct{}{}
		defer func() { <-workerPool }()
		extractingStartMessage := entities.NewMessage(entities.TargetVideoAPIService, entities.StartProcessingMessage, nil)
		err := converVideo.videoProcessorMessaging.Publish(extractingStartMessage)
		if err != nil {
			//TODO remover panic
			panic(err)
		}
		sourceFile := entities.NewFile(ConvertVideoInput.VideoId, ConvertVideoInput.VideoName, ConvertVideoInput.VideoPath)
		newFile, err := converVideo.videoProcessorStorage.Download(sourceFile)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)
			return
		}
		defer os.Remove(filepath.Join(newFile.Path, newFile.Name))

		newFolder, err := converVideo.videoProcessorConverter.ConvertToImages(newFile)
		defer os.RemoveAll(newFolder.Path)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)

			return
		}
		compressedFile, err := converVideo.videoProcessorCompressor.Compress(newFolder)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)

			return
		}
		uploadURL, err := converVideo.videoProcessorStorage.Upload(compressedFile)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)

			return
		}
		defer os.Remove(compressedFile.Path + "")
		extractingSuccessMessage := entities.NewMessage(
			entities.TargetVideoAPIService,
			entities.ExtractSuccessMessage,
			entities.ExtractSuccessPayload{
				VideoSnapshotsUrl: uploadURL,
				VideoId:           ConvertVideoInput.VideoId,
			})

		converVideo.videoProcessorMessaging.Publish(extractingSuccessMessage)

	}()

	return ConvertVideoOutput{}, nil
}

func (converVideo *ConvertVideoUsecase) SendErrorMessage(err error, ConvertVideoInput ConvertVideoInput) {

	processingErrorMessage := entities.NewMessage(
		entities.TargetVideoAPIService,
		entities.ExtractErrorMessage,
		entities.ExtractErrorPayload{
			VideoId:          ConvertVideoInput.VideoId,
			ErrorMessage:     err.Error(),
			ErrorDescription: err.Error(),
		})
	sendErrorMessage := entities.NewMessage(
		entities.TargetEmailService,
		entities.SendErrorMessage,
		entities.ExtractSendErrorPayload{
			UserID:            ConvertVideoInput.UserId,
			VideoUrl:          ConvertVideoInput.VideoPath,
			VideoSnapshotsUrl: ConvertVideoInput.VideoPath,
			VideoDescription:  ConvertVideoInput.VideoDescription,
			ErrorMessage:      err.Error(),
			ErrorDescription:  err.Error(),
		})
	converVideo.videoProcessorMessaging.Publish(processingErrorMessage)
	converVideo.videoProcessorMessaging.Publish(sendErrorMessage)

}
