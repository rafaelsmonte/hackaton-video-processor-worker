package mappers

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/sqs/dto"
)

func ProcessVideoInput(u dto.VideoProcessRequest) usecases.ConvertVideoInput {
	return usecases.ConvertVideoInput{
		VideoId:   u.VideoId,
		UserId:    u.UserId,
		VideoName: u.VideoName,
	}
}

func ProcessVideoResponse(u usecases.ConvertVideoOutput) dto.VideoProcessResponse {
	return dto.VideoProcessResponse{
		Message: "Video processed successfully at " + u.VideoUrl,
	}
}
