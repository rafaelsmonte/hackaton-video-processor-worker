package sqs

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/FFMPEG"
	"hackaton-video-processor-worker/internal/infra/s3"
	"hackaton-video-processor-worker/internal/infra/sns"
	"hackaton-video-processor-worker/internal/infra/sqs/handlers"
	"hackaton-video-processor-worker/internal/infra/zip"
	"log"
)

func ConfigProcessingHandlers() *AppHandlers {
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mqRepository, err := sns.NewSNS()
	if err != nil {
		log.Fatalln("Error connecting to SNS", err)
	}
	storageRepository, err := s3.NewS3()
	if err != nil {
		log.Fatalln("Error connecting to SNS", err)
	}
	zipRepository := zip.NewZIP()
	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, mqRepository, storageRepository, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)
	return &AppHandlers{
		VideoProcessorHandler: videoProcessorHandler,
	}
}
func ConfigDLQHandlers() *AppHandlers {
	mqRepository, err := sns.NewSNS()
	if err != nil {
		log.Fatalln("Error connecting to SNS", err)
	}

	videoUsecase := usecases.NewGenericErrorUsecase(mqRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)
	return &AppHandlers{
		VideoProcessorHandler: videoProcessorHandler,
	}
}
