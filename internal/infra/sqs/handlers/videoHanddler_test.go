package handlers_test

import (
	"encoding/json"
	"errors"
	"testing"

	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/sqs/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConvertVideoUsecase struct {
	mock.Mock
}

func (m *MockConvertVideoUsecase) Execute(input usecases.ConvertVideoInput) (usecases.ConvertVideoOutput, error) {
	args := m.Called(input)
	return args.Get(0).(usecases.ConvertVideoOutput), args.Error(1)
}

func TestHandleMessage_InvalidJson(t *testing.T) {
	mockUsecase := new(MockConvertVideoUsecase)
	handler := handlers.NewVideoHandler(mockUsecase)
	invalidBody := "{ invalid json }"

	err := handler.HandleMessage(&invalidBody)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid message format")
}

func TestHandleMessage_UnsupportedMessageType(t *testing.T) {
	mockUsecase := new(MockConvertVideoUsecase)
	handler := handlers.NewVideoHandler(mockUsecase)

	message := struct {
		Sender  string `json:"sender"`
		Target  string `json:"target"`
		Type    string `json:"type"`
		Payload struct {
			VideoUrl string `json:"videoUrl"`
			VideoId  string `json:"videoId"`
		} `json:"payload"`
	}{
		Type: "UNKNOWN_TYPE",
	}

	bodyBytes, _ := json.Marshal(message)
	body := string(bodyBytes)

	// Act
	err := handler.HandleMessage(&body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported message type")
}

func TestHandleMessage_SuccessfulProcessing(t *testing.T) {
	mockUsecase := new(MockConvertVideoUsecase)
	handler := handlers.NewVideoHandler(mockUsecase)

	message := struct {
		Sender  string `json:"sender"`
		Target  string `json:"target"`
		Type    string `json:"type"`
		Payload struct {
			VideoUrl string `json:"videoUrl"`
			VideoId  string `json:"videoId"`
		} `json:"payload"`
	}{
		Type: "MSG_EXTRACT_SNAPSHOT",
		Payload: struct {
			VideoUrl string `json:"videoUrl"`
			VideoId  string `json:"videoId"`
		}{
			VideoUrl: "http://example.com/video.mp4",
			VideoId:  "1234",
		},
	}

	bodyBytes, _ := json.Marshal(message)
	body := string(bodyBytes)

	mockOutput := usecases.ConvertVideoOutput{
		VideoUrl: "http://example.com/video.mp4",
	}

	mockUsecase.On("Execute", mock.Anything).Return(mockOutput, nil)

	err := handler.HandleMessage(&body)

	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}

func TestHandleMessage_FailedProcessing(t *testing.T) {
	mockUsecase := new(MockConvertVideoUsecase)
	handler := handlers.NewVideoHandler(mockUsecase)

	message := struct {
		Sender  string `json:"sender"`
		Target  string `json:"target"`
		Type    string `json:"type"`
		Payload struct {
			VideoUrl string `json:"videoUrl"`
			VideoId  string `json:"videoId"`
		} `json:"payload"`
	}{
		Type: "MSG_EXTRACT_SNAPSHOT",
		Payload: struct {
			VideoUrl string `json:"videoUrl"`
			VideoId  string `json:"videoId"`
		}{
			VideoUrl: "http://example.com/video.mp4",
			VideoId:  "1234",
		},
	}

	bodyBytes, _ := json.Marshal(message)
	body := string(bodyBytes)

	mockUsecase.On("Execute", mock.Anything).Return(usecases.ConvertVideoOutput{}, errors.New("processing error"))

	err := handler.HandleMessage(&body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "video processing error")
	mockUsecase.AssertExpectations(t)
}
