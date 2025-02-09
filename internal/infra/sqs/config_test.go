package sqs

import (
	"hackaton-video-processor-worker/internal/domain/entities"
	"hackaton-video-processor-worker/internal/domain/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFFMPEG é um mock para FFMPEG
type MockFFMPEG struct {
	mock.Mock
}

func (m *MockFFMPEG) ProcessVideo(input string) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

// MockSNS é um mock para SNS
type MockSNS struct {
	mock.Mock
}

func (m *MockSNS) Publish(message string) error {
	args := m.Called(message)
	return args.Error(0)
}

// MockS3 é um mock para S3
type MockS3 struct {
	mock.Mock
}

func (m *MockS3) UploadFile(file string) error {
	args := m.Called(file)
	return args.Error(0)
}

// MockZIP é um mock para ZIP
type MockZIP struct {
	mock.Mock
}

func (m *MockZIP) Compress(files []string) (string, error) {
	args := m.Called(files)
	return args.String(0), args.Error(1)
}

// MockConvertVideoUsecase é um mock para usecases.ConvertVideoUsecase
type MockConvertVideoUsecase struct {
	mock.Mock
}

func (m *MockConvertVideoUsecase) Execute(input string) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

// MockGenericErrorUsecase é um mock para usecases.GenericErrorUsecase
type MockGenericErrorUsecase struct {
	mock.Mock
}

func (m *MockGenericErrorUsecase) Execute(input string) error {
	args := m.Called(input)
	return args.Error(0)
}

type MockCompressor struct {
	mock.Mock
}

func (m *MockCompressor) Compress(folder entities.Folder) (entities.File, error) {
	args := m.Called(folder)
	return args.Get(0).(entities.File), args.Error(1)
}

type MockConverter struct {
	mock.Mock
}

func (m *MockConverter) ConvertToImages(file entities.File) (entities.Folder, error) {
	args := m.Called(file)
	return args.Get(0).(entities.Folder), args.Error(1)
}

type MockMessaging struct {
	mock.Mock
}

func (m *MockMessaging) Publish(message entities.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Download(file entities.File) (entities.File, error) {
	args := m.Called(file)
	return args.Get(0).(entities.File), args.Error(1)
}

func (m *MockStorage) Upload(file entities.File) (string, error) {
	args := m.Called(file)
	return args.String(0), args.Error(1)
}

func TestConfigProcessingHandlers(t *testing.T) {
	// Cria os mocks
	mockConverter := new(MockConverter)
	mockMessaging := new(MockMessaging)
	mockStorage := new(MockStorage)
	mockCompressor := new(MockCompressor)

	// Configura as expectativas dos mocks
	mockConverter.On("ConvertToImages", entities.File{Path: "input.mp4"}).Return(entities.Folder{Path: "frames"}, nil)
	mockMessaging.On("Publish", entities.Message(entities.Message{})).Return(nil)
	mockStorage.On("Download", entities.File{Path: "input.mp4"}).Return(entities.File{Path: "input.mp4"}, nil)
	mockStorage.On("Upload", entities.File{Path: "output.mp4"}).Return("output.mp4", nil)
	mockCompressor.On("Compress", entities.Folder{Path: "frames"}).Return(entities.File{Path: "output.zip"}, nil)

	// Injeta os mocks no use case
	_ = usecases.NewConvertVideoUsecase(mockConverter, mockMessaging, mockStorage, mockCompressor)

	// Chama a função a ser testada
	handlers := ConfigProcessingHandlers()

	// Verifica se os handlers foram configurados corretamente
	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.VideoProcessorHandler)

}

func TestConfigDLQHandlers(t *testing.T) {
	// Cria os mocks
	mockMessaging := new(MockMessaging)

	// Configura as expectativas dos mocks
	mockMessaging.On("Publish", entities.Message{}).Return(nil)

	// Injeta os mocks no use case
	_ = usecases.NewGenericErrorUsecase(mockMessaging)

	// Chama a função a ser testada
	handlers := ConfigDLQHandlers()

	// Verifica se os handlers foram configurados corretamente
	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.VideoProcessorHandler)

}
