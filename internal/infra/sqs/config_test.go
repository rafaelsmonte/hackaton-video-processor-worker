package sqs

import (
	"errors"
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/FFMPEG"
	"hackaton-video-processor-worker/internal/infra/s3"
	"hackaton-video-processor-worker/internal/infra/sns"

	//"hackaton-video-processor-worker/internal/infra/sqs"
	"hackaton-video-processor-worker/internal/infra/sqs/handlers"
	"hackaton-video-processor-worker/internal/infra/zip"
	"testing"
)

// Mock para SNS
type mockSNS struct {
	err error
}

func (m *mockSNS) NewSNS() (*sns.SNS, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &sns.SNS{}, nil
}

// Mock para S3
type mockS3 struct {
	err error
}

func (m *mockS3) NewS3() (*s3.S3, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &s3.S3{}, nil
}

func TestConfigHandlers_Success(t *testing.T) {
	// Criando instâncias dos mocks
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mockSNS := &mockSNS{}
	mockS3 := &mockS3{}
	zipRepository := zip.NewZIP()

	// Simulando instâncias reais sem erro
	mqRepository, err := mockSNS.NewSNS()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	storageRepository, err := mockS3.NewS3()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, mqRepository, storageRepository, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)

	appHandlers := &AppHandlers{
		VideoProcessorHandler: videoProcessorHandler,
	}

	if appHandlers.VideoProcessorHandler == nil {
		t.Fatal("Expected videoProcessorHandler to be initialized, but got nil")
	}
}

func TestConfigHandlers_FailOnSNS(t *testing.T) {
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mockSNS := &mockSNS{err: errors.New("SNS connection failed")}
	mockS3 := &mockS3{}
	zipRepository := zip.NewZIP()

	_, err := mockSNS.NewSNS()
	if err == nil {
		t.Fatal("Expected error for SNS connection, but got nil")
	}

	storageRepository, err := mockS3.NewS3()
	if err != nil {
		t.Fatalf("Expected no error for S3, but got: %v", err)
	}

	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, nil, storageRepository, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)

	appHandlers := &AppHandlers{
		VideoProcessorHandler: videoProcessorHandler,
	}

	if appHandlers.VideoProcessorHandler == nil {
		t.Fatal("Expected videoProcessorHandler to be initialized, but got nil")
	}
}

func TestConfigHandlers_FailOnS3(t *testing.T) {
	videoProcessorRepository := FFMPEG.NewFFMPEG()
	mockSNS := &mockSNS{}
	mockS3 := &mockS3{err: errors.New("S3 connection failed")}
	zipRepository := zip.NewZIP()

	mqRepository, err := mockSNS.NewSNS()
	if err != nil {
		t.Fatalf("Expected no error for SNS, but got: %v", err)
	}

	_, err = mockS3.NewS3()
	if err == nil {
		t.Fatal("Expected error for S3 connection, but got nil")
	}

	videoUsecase := usecases.NewConvertVideoUsecase(videoProcessorRepository, mqRepository, nil, zipRepository)
	videoProcessorHandler := handlers.NewVideoHandler(&videoUsecase)

	appHandlers := &AppHandlers{
		VideoProcessorHandler: videoProcessorHandler,
	}

	if appHandlers.VideoProcessorHandler == nil {
		t.Fatal("Expected videoProcessorHandler to be initialized, but got nil")
	}
}
