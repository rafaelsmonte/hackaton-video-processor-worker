package usecases_test

import (
	"hackaton-video-processor-worker/internal/domain/entities"
	"hackaton-video-processor-worker/internal/domain/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVideoProcessorMessaging struct {
	mock.Mock
}

func (m *MockVideoProcessorMessaging) Publish(message entities.Message) error {
	m.Called(message)
	return nil
}

func TestGenericErrorUsecase_Execute(t *testing.T) {
	mockMessaging := new(MockVideoProcessorMessaging)
	usecase := usecases.NewGenericErrorUsecase(mockMessaging)

	input := usecases.ConvertVideoInput{
		VideoId: "12345",
		UserId:  "user-1",
	}

	expectedMessage := entities.NewMessage(
		entities.TargetVideoSQSService,
		entities.ExtractErrorMessage,
		entities.ExtractErrorPayload{
			VideoId:          input.VideoId,
			UserId:           input.UserId,
			ErrorMessage:     "GENERIC_ERROR",
			ErrorDescription: "Gerneric Error",
		},
	)

	mockMessaging.On("Publish", expectedMessage).Return()

	output, err := usecase.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, usecases.ConvertVideoOutput{}, output)
	mockMessaging.AssertExpectations(t)
}
