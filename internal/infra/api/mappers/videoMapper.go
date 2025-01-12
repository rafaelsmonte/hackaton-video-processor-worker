package mappers

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/api/dto"
)

func ProcessVideoInput(u dto.VideoProcessRequest) usecases.ConvertVideoInput {
	return usecases.ConvertVideoInput{
		VideoName: u.Name,
		VideoPath: u.Path,
		VideoId:   u.Id,
	}
}

func ProcessVideoResponse(u usecases.ConvertVideoOutput) dto.VideoProcessResponse {
	return dto.VideoProcessResponse{
		Message: u.VideoPath + "Video processando",
	}
}
