package mappers

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/sqs/dto"
)

func ProcessVideoInput(u dto.VideoProcessRequest) usecases.ConvertVideoInput {
	return usecases.ConvertVideoInput{
		VideoUrl: u.VideoUrl,
		VideoId:  u.VideoId,
	}
}

func ProcessVideoResponse(u usecases.ConvertVideoOutput) dto.VideoProcessResponse {
	return dto.VideoProcessResponse{
		Message: "Video processed successfully at " + u.VideoUrl,
	}
}
