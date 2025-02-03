package usecases

import (
	"errors"
	"hackaton-video-processor-worker/internal/domain/entities"
	"hackaton-video-processor-worker/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestConvertVideoUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConverter := mocks.NewMockIVideoProcessorConverter(ctrl)
	mockMessaging := mocks.NewMockIVideoProcessorMessaging(ctrl)
	mockStorage := mocks.NewMockIVideoProcessorStorage(ctrl)
	mockCompressor := mocks.NewMockIVideoProcessorCompressor(ctrl)

	usecase := ConvertVideoUsecase{
		videoProcessorConverter:  mockConverter,
		videoProcessorMessaging:  mockMessaging,
		videoProcessorStorage:    mockStorage,
		videoProcessorCompressor: mockCompressor,
	}

	input := ConvertVideoInput{
		VideoName:        "test_video",
		VideoUrl:         "http://example.com/video.mp4",
		VideoId:          "12345",
		UserId:           "67890",
		VideoDescription: "A test video",
	}

	mockMessaging.EXPECT().Publish(gomock.Any()).Return(nil).Times(1)
	mockStorage.EXPECT().Download(gomock.Any()).Return(entities.File{Id: "12345", Path: "/tmp/12345"}, nil).Times(1)
	mockConverter.EXPECT().ConvertToImages(gomock.Any()).Return(entities.Folder{Name: "12345", Path: "/tmp/12345_images"}, nil).Times(1)
	mockCompressor.EXPECT().Compress(gomock.Any()).Return(entities.File{Id: "12345_compressed", Path: "/tmp/12345_compressed"}, nil).Times(1)
	mockStorage.EXPECT().Upload(gomock.Any()).Return("http://example.com/uploaded_video.mp4", nil).Times(1)
	mockMessaging.EXPECT().Publish(gomock.Any()).Return(nil).Times(1)

	wg.Add(1)
	output, err := usecase.Execute(input)
	wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, ConvertVideoOutput{}, output)
}

func TestConvertVideoUsecase_Execute_ErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConverter := mocks.NewMockIVideoProcessorConverter(ctrl)
	mockMessaging := mocks.NewMockIVideoProcessorMessaging(ctrl)
	mockStorage := mocks.NewMockIVideoProcessorStorage(ctrl)
	mockCompressor := mocks.NewMockIVideoProcessorCompressor(ctrl)

	usecase := ConvertVideoUsecase{
		videoProcessorConverter:  mockConverter,
		videoProcessorMessaging:  mockMessaging,
		videoProcessorStorage:    mockStorage,
		videoProcessorCompressor: mockCompressor,
	}

	input := ConvertVideoInput{
		VideoName:        "test_video",
		VideoUrl:         "http://example.com/video.mp4",
		VideoId:          "12345",
		UserId:           "67890",
		VideoDescription: "A test video",
	}

	mockMessaging.EXPECT().Publish(gomock.Any()).Return(nil).Times(1)
	mockStorage.EXPECT().Download(gomock.Any()).Return(entities.File{}, errors.New("download error")).Times(1)
	mockMessaging.EXPECT().Publish(gomock.Any()).Return(nil).Times(2) // Espera que duas mensagens de erro sejam publicadas

	wg.Add(1)
	output, err := usecase.Execute(input)
	wg.Wait()

	assert.NoError(t, err)
	assert.Equal(t, ConvertVideoOutput{}, output)
}

func TestConvertVideoUsecase_SendErrorMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessaging := mocks.NewMockIVideoProcessorMessaging(ctrl)

	usecase := ConvertVideoUsecase{
		videoProcessorMessaging: mockMessaging,
	}

	input := ConvertVideoInput{
		VideoName:        "test_video",
		VideoUrl:         "http://example.com/video.mp4",
		VideoId:          "12345",
		UserId:           "67890",
		VideoDescription: "A test video",
	}

	err := errors.New("test error")

	mockMessaging.EXPECT().Publish(gomock.Any()).Return(nil).Times(2)

	usecase.SendErrorMessage(err, input)
}
