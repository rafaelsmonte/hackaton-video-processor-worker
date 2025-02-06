package sqs

import (
	"context"
	"errors"
	"hackaton-video-processor-worker/internal/domain/usecases"

	//sqsInteance "hackaton-video-processor-worker/internal/infra/sqs"
	"hackaton-video-processor-worker/internal/infra/sqs/handlers"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQSClient struct {
	mock.Mock
}
type MockSQSService struct {
	mock.Mock
}

func (m *MockSQSService) StartConsuming(ctx context.Context) {
	m.Called(ctx)
}
func (m *MockSQSClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

func (m *MockSQSClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	args := m.Called(input)
	return &sqs.DeleteMessageOutput{}, args.Error(1)
}

type MockVideoHandler struct {
	handlers.VideoHandler
	mock.Mock
}

func (m *MockVideoHandler) HandleMessage(msg *string) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockVideoHandler) Execute(ConvertVideoInput usecases.ConvertVideoInput) (usecases.ConvertVideoOutput, error) {
	args := m.Called(ConvertVideoInput)
	return usecases.ConvertVideoOutput{}, args.Error(0)
}

type MockVideoUseCaseHanddler struct {
	mock.Mock
}

func TestReceiveMessages_Success(t *testing.T) {
	mockSQS := new(MockSQSClient)
	mockHandler := new(MockVideoHandler)

	sqsService := &SQSService{
		SqsClient: mockSQS,
		QueueURL:  "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
		Handler:   &AppHandlers{VideoProcessorHandler: &mockHandler.VideoHandler},
	}

	mockMessages := &sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{Body: aws.String("test-message"), ReceiptHandle: aws.String("handle-123")},
		},
	}

	// Set expectations for the mockSQS methods
	mockSQS.On("ReceiveMessage", mock.Anything).Return(mockMessages, nil)
	mockHandler.On("HandleMessage", mock.Anything).Return(nil)
	mockSQS.On("DeleteMessage", mock.MatchedBy(func(input *sqs.DeleteMessageInput) bool {
		return *input.ReceiptHandle == "handle-123"
	})).Return(&sqs.DeleteMessageOutput{}, nil)

	messages, err := sqsService.ReceiveMessages()

	assert.NoError(t, err)
	assert.Len(t, messages, 1)

}

func TestProcessMessage_Failure_HandleMessage(t *testing.T) {
	mockSQS := new(MockSQSClient)
	mockHandler := new(MockVideoHandler)

	sqsService := &SQSService{
		SqsClient: mockSQS,
		QueueURL:  "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
		Handler:   &AppHandlers{VideoProcessorHandler: &mockHandler.VideoHandler},
	}

	mockHandler.On("HandleMessage", mock.Anything).Return(errors.New("processing error"))

	err := sqsService.ProcessMessage(&sqs.Message{Body: aws.String("test-message"), ReceiptHandle: aws.String("handle-123")})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to handle message")
}
func TestDeleteMessage(t *testing.T) {
	mockSQS := new(MockSQSClient)

	sqsService := &SQSService{
		SqsClient: mockSQS,
		QueueURL:  "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}

	receiptHandle := aws.String("handle-123")

	mockSQS.On("DeleteMessage", mock.Anything).Return(&sqs.DeleteMessageOutput{}, nil)

	err := sqsService.DeleteMessage(receiptHandle)

	assert.NoError(t, err)

	mockSQS.AssertExpectations(t)
}
