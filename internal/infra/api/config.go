package api

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/FFMPEG"
	"hackaton-video-processor-worker/internal/infra/api/handlers"
	"hackaton-video-processor-worker/internal/infra/fakeS3"
	"hackaton-video-processor-worker/internal/infra/sns"
	"hackaton-video-processor-worker/internal/infra/zip"
	"log"
)

func configHandlers() *AppHandlers {
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mqRepository, err := sns.NewSNS()
	if err != nil {
		log.Fatalln("Error connecting to SNS", err)
	}
	storageRepository, err := fakeS3.NewS3()
	if err != nil {
		log.Fatalln("Error connecting to SNS", err)
	}
	zipRepository := zip.NewZIP()
	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, mqRepository, storageRepository, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)

	return &AppHandlers{
		videoProcessorHandler: videoProcessorHandler,
	}
}
