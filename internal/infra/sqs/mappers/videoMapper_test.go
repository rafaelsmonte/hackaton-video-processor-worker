package mappers

import (
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/sqs/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessVideoInput(t *testing.T) {
	input := dto.VideoProcessRequest{
		VideoId: "12345",
	}

	expected := usecases.ConvertVideoInput{
		VideoId: "12345",
	}

	result := ProcessVideoInput(input)

	assert.Equal(t, expected, result)
}

func TestProcessVideoResponse(t *testing.T) {
	output := usecases.ConvertVideoOutput{
		VideoUrl: "http://example.com/processed.mp4",
	}

	expected := dto.VideoProcessResponse{
		Message: "Video processed successfully at http://example.com/processed.mp4",
	}

	result := ProcessVideoResponse(output)

	assert.Equal(t, expected, result)
}
