package sqs

import (
	"context"
	"errors"
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/sqs/handlers"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSQSClient implements SQSClient interface for testing
type MockSQSClient struct {
	mock.Mock
}

func (m *MockSQSClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(input)
	if output := args.Get(0); output != nil {
		return output.(*sqs.ReceiveMessageOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSQSClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	args := m.Called(input)
	if output := args.Get(0); output != nil {
		return output.(*sqs.DeleteMessageOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockVideoUsecase for testing
type MockVideoUsecase struct {
	mock.Mock
}

func (m *MockVideoUsecase) Execute(input usecases.ConvertVideoInput) (usecases.ConvertVideoOutput, error) {
	args := m.Called(input)
	return args.Get(0).(usecases.ConvertVideoOutput), args.Error(1)
}

func TestNewSQSService(t *testing.T) {
	// Arrange
	region := "us-east-1"
	queueURL := "https://sqs.example.com/queue"
	mockUsecase := new(MockVideoUsecase)
	handler := &AppHandlers{
		VideoProcessorHandler: handlers.NewVideoHandler(mockUsecase),
	}

	// Act
	service := NewSQSService(region, queueURL, handler)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, queueURL, service.QueueURL)
	assert.NotNil(t, service.SqsClient)
	assert.Equal(t, handler, service.Handler)
}

func TestReceiveMessages(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockSQSClient)
		expectedMsgs  []*sqs.Message
		expectedError error
	}{
		{
			name: "successful_receive",
			setupMock: func(m *MockSQSClient) {
				m.On("ReceiveMessage", &sqs.ReceiveMessageInput{
					QueueUrl:            aws.String("test-queue"),
					MaxNumberOfMessages: aws.Int64(10),
					WaitTimeSeconds:     aws.Int64(1),
				}).Return(&sqs.ReceiveMessageOutput{
					Messages: []*sqs.Message{
						{
							MessageId:     aws.String("msg1"),
							Body:          aws.String(`{"type":"MSG_EXTRACT_SNAPSHOT","payload":{"userId":"123","videoId":"456","videoName":"test.mp4"}}`),
							ReceiptHandle: aws.String("receipt1"),
						},
					},
				}, nil)
			},
			expectedMsgs: []*sqs.Message{
				{
					MessageId:     aws.String("msg1"),
					Body:          aws.String(`{"type":"MSG_EXTRACT_SNAPSHOT","payload":{"userId":"123","videoId":"456","videoName":"test.mp4"}}`),
					ReceiptHandle: aws.String("receipt1"),
				},
			},
			expectedError: nil,
		},
		{
			name: "receive_error",
			setupMock: func(m *MockSQSClient) {
				m.On("ReceiveMessage", mock.Anything).Return(nil, errors.New("receive error"))
			},
			expectedMsgs:  nil,
			expectedError: errors.New("receive error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := new(MockSQSClient)
			tt.setupMock(mockClient)

			mockUsecase := new(MockVideoUsecase)
			service := &SQSService{
				SqsClient: mockClient,
				QueueURL:  "test-queue",
				Handler: &AppHandlers{
					VideoProcessorHandler: handlers.NewVideoHandler(mockUsecase),
				},
			}

			// Act
			messages, err := service.ReceiveMessages()

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsgs, messages)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestProcessMessage(t *testing.T) {
	validMessageBody := `{"type":"MSG_EXTRACT_SNAPSHOT","payload":{"userId":"123","videoId":"456","videoName":"test.mp4"}}`

	tests := []struct {
		name          string
		message       *sqs.Message
		setupMocks    func(*MockSQSClient, *MockVideoUsecase)
		expectedError error
	}{
		{
			name: "successful_processing",
			message: &sqs.Message{
				MessageId:     aws.String("msg1"),
				Body:          aws.String(validMessageBody),
				ReceiptHandle: aws.String("receipt1"),
			},
			setupMocks: func(sqsMock *MockSQSClient, usecaseMock *MockVideoUsecase) {
				usecaseMock.On("Execute", mock.Anything).Return(usecases.ConvertVideoOutput{
					VideoUrl: "Video_URL",
				}, nil)
				sqsMock.On("DeleteMessage", &sqs.DeleteMessageInput{
					QueueUrl:      aws.String("test-queue"),
					ReceiptHandle: aws.String("receipt1"),
				}).Return(&sqs.DeleteMessageOutput{}, nil)
			},
			expectedError: nil,
		},
		{
			name: "usecase_error",
			message: &sqs.Message{
				MessageId:     aws.String("msg1"),
				Body:          aws.String(validMessageBody),
				ReceiptHandle: aws.String("receipt1"),
			},
			setupMocks: func(sqsMock *MockSQSClient, usecaseMock *MockVideoUsecase) {
				usecaseMock.On("Execute", mock.Anything).Return(usecases.ConvertVideoOutput{}, errors.New("usecase error"))
			},
			expectedError: errors.New("failed to handle message: video processing error: usecase error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := new(MockSQSClient)
			mockUsecase := new(MockVideoUsecase)
			tt.setupMocks(mockClient, mockUsecase)

			service := &SQSService{
				SqsClient: mockClient,
				QueueURL:  "test-queue",
				Handler: &AppHandlers{
					VideoProcessorHandler: handlers.NewVideoHandler(mockUsecase),
				},
			}

			// Act
			err := service.ProcessMessage(tt.message)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestStartConsuming(t *testing.T) {
	// Arrange
	mockClient := new(MockSQSClient)
	mockUsecase := new(MockVideoUsecase)
	validMessageBody := `{"type":"MSG_EXTRACT_SNAPSHOT","payload":{"userId":"123","videoId":"456","videoName":"test.mp4"}}`

	mockClient.On("ReceiveMessage", mock.Anything).Return(&sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{
				MessageId:     aws.String("msg1"),
				Body:          aws.String(validMessageBody),
				ReceiptHandle: aws.String("receipt1"),
			},
		},
	}, nil).Maybe()

	mockUsecase.On("Execute", mock.Anything).Return(usecases.ConvertVideoOutput{
		VideoUrl: "Video_URL",
	}, nil).Maybe()

	mockClient.On("DeleteMessage", mock.Anything).Return(&sqs.DeleteMessageOutput{}, nil).Maybe()

	service := &SQSService{
		SqsClient: mockClient,
		QueueURL:  "test-queue",
		Handler: &AppHandlers{
			VideoProcessorHandler: handlers.NewVideoHandler(mockUsecase),
		},
	}

	// Act
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start consuming in a goroutine
	go service.StartConsuming(ctx)

	// Wait for context cancellation
	<-ctx.Done()

	// Assert
	mockClient.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}
