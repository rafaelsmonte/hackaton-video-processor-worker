package api

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/FFMPEG"
	"hackaton-video-processor-worker/internal/infra/api/handlers"
	"hackaton-video-processor-worker/internal/infra/fakeMQ"
	"hackaton-video-processor-worker/internal/infra/s3"
	"hackaton-video-processor-worker/internal/infra/zip"
)

func configHandlers() *AppHandlers {
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mqRepository := fakeMQ.NewFakeMQ()
	storageRepository := s3.NewS3()
	zipRepository := zip.NewZIP()
	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, mqRepository, storageRepository, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)

	return &AppHandlers{
		videoProcessorHandler: videoProcessorHandler,
	}
}
