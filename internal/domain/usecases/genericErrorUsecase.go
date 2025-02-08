package usecases

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"log"
)

type IGenericErrorUsecase interface {
	Execute(ConvertVideoInput ConvertVideoInput) (ConvertVideoOutput, error)
}
type GenericErrorUsecase struct {
	videoProcessorMessaging adapters.IVideoProcessorMessaging
}

func (genericError *GenericErrorUsecase) Execute(ConvertVideoInput ConvertVideoInput) (ConvertVideoOutput, error) {

	genericErrorMessage := entities.NewMessage(
		entities.TargetVideoSQSService,
		entities.ExtractErrorMessage,
		entities.ExtractErrorPayload{
			VideoId:          ConvertVideoInput.VideoId,
			UserId:           ConvertVideoInput.UserId,
			ErrorMessage:     "GENERIC_ERROR",
			ErrorDescription: "Gerneric Error",
		})
	log.Println("Generic Error <<<", genericErrorMessage)

	genericError.videoProcessorMessaging.Publish(genericErrorMessage)

	return ConvertVideoOutput{}, nil
}

type GenericErrorInput struct {
	VideoName        string
	VideoUrl         string
	VideoId          string
	UserId           string
	VideoDescription string
}

type GenericErrorOutput struct {
	VideoUrl string
}

func NewGenericErrorUsecase(
	videoProcessorMessaging adapters.IVideoProcessorMessaging,
) GenericErrorUsecase {
	return GenericErrorUsecase{
		videoProcessorMessaging: videoProcessorMessaging,
	}
}
