package api

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/FFMPEG"
	"hackaton-video-processor-worker/internal/infra/api/handlers"
	"hackaton-video-processor-worker/internal/infra/fakeMQ"
	"hackaton-video-processor-worker/internal/infra/fakeStorage"
	"hackaton-video-processor-worker/internal/infra/fakeZIP"
)

func configHandlers() *AppHandlers {
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mqRepository := fakeMQ.NewFakeMQ()
	storageRepository := fakeStorage.NewFakeStorage()
	zipRepository := fakeZIP.NewFakeZIP()
	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, mqRepository, storageRepository, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)

	return &AppHandlers{
		videoProcessorHandler: videoProcessorHandler,
	}
}
