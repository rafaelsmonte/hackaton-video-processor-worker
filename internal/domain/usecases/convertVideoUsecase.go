package usecases

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

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

var workerPool = make(chan struct{}, getWorkerPoolSize())

func getWorkerPoolSize() int {

	poolSizeStr := os.Getenv("WORKER_POOL_SIZE")
	if poolSizeStr == "" {
		return 5
	}
	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		return 5
	}
	return poolSize
}
func (converVideo *ConvertVideoUsecase) Execute(ConvertVideoInput ConvertVideoInput) (ConvertVideoOutput, error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
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
			return
		}
		sourceFile := entities.NewFile(ConvertVideoInput.VideoId, ConvertVideoInput.VideoUrl, ConvertVideoInput.UserId, ConvertVideoInput.VideoName, nil)
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
		defer os.Remove(compressedFile.Name)
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

	converVideo.videoProcessorMessaging.Publish(processingErrorMessage)

}
