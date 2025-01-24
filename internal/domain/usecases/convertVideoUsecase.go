package usecases

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
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
	VideoUrl         string
	VideoId          string
	UserId           string
	VideoDescription string
}

type ConvertVideoOutput struct {
	VideoUrl string
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

		extractingStartMessage := entities.NewMessage(
			entities.TargetVideoSQSService,
			entities.StartProcessingMessage,
			entities.StartProcessingPayload{
				VideoId: ConvertVideoInput.VideoId,
			})
		err := converVideo.videoProcessorMessaging.Publish(extractingStartMessage)
		if err != nil {
			//TODO remover panic
			panic(err)
		}
		sourceFile := entities.NewFile(ConvertVideoInput.VideoId, ConvertVideoInput.VideoUrl, nil)
		newFile, err := converVideo.videoProcessorStorage.Download(sourceFile)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)
			return
		}

		newFolder, err := converVideo.videoProcessorConverter.ConvertToImages(newFile)
		defer os.RemoveAll(newFolder.Path)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)

			return
		}
		compressedFile, err := converVideo.videoProcessorCompressor.Compress(newFolder)
		defer os.Remove(compressedFile.Id)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)
			return
		}
		uploadURL, err := converVideo.videoProcessorStorage.Upload(compressedFile)
		if err != nil {
			converVideo.SendErrorMessage(err, ConvertVideoInput)

			return
		}
		extractingSuccessMessage := entities.NewMessage(
			entities.TargetVideoSQSService,
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
		entities.TargetVideoSQSService,
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
			VideoUrl:          ConvertVideoInput.VideoUrl,
			VideoSnapshotsUrl: ConvertVideoInput.VideoUrl,
			ErrorMessage:      err.Error(),
			ErrorDescription:  err.Error(),
		})
	converVideo.videoProcessorMessaging.Publish(processingErrorMessage)
	converVideo.videoProcessorMessaging.Publish(sendErrorMessage)

}
