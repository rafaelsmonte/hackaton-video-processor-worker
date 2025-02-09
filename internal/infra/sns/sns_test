package sns

import (
	"context"
	"encoding/json"
	"errors"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSNSClient struct {
	mock.Mock
}

func (m *mockSNSClient) Publish(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, input)
	return &sns.PublishOutput{}, args.Error(1)
}

func TestPublish_Success(t *testing.T) {
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")

	mockClient := new(mockSNSClient)
	snsInstance := &SNS{Client: mockClient}

	message := entities.Message{
		Type:    "TestType",
		Payload: map[string]string{"key": "value"},
	}

	messageBody, _ := json.Marshal(message)
	mockClient.On("Publish", mock.Anything, &sns.PublishInput{
		TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
		Message:  aws.String(string(messageBody)),
	}).Return(&sns.PublishOutput{}, nil)

	err := snsInstance.Publish(message)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestPublish_FailSerialization(t *testing.T) {
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")

	snsInstance := &SNS{}

	message := entities.Message{
		Type:    "TestType",
		Payload: make(chan int),
	}

	err := snsInstance.Publish(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to serialize message")
}

func TestPublish_FailToPublish(t *testing.T) {
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")

	mockClient := new(mockSNSClient)
	snsInstance := &SNS{Client: mockClient}

	message := entities.Message{
		Type:    "TestType",
		Payload: map[string]string{"key": "value"},
	}

	messageBody, _ := json.Marshal(message)
	mockClient.On("Publish", mock.Anything, &sns.PublishInput{
		TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
		Message:  aws.String(string(messageBody)),
	}).Return(&sns.PublishOutput{}, errors.New("publish error"))

	err := snsInstance.Publish(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message")
	mockClient.AssertExpectations(t)
}

func TestPublish_MissingTopicArn(t *testing.T) {
	os.Unsetenv("SNS_TOPIC_ARN")

	mockClient := new(mockSNSClient)
	snsInstance := &SNS{Client: mockClient}

	message := entities.Message{
		Type:    "TestType",
		Payload: map[string]string{"key": "value"},
	}

	err := snsInstance.Publish(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SNS_TOPIC_ARN is not set")
}

func TestPublish_LoggingSerializationError(t *testing.T) {
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")

	mockClient := new(mockSNSClient)
	snsInstance := &SNS{Client: mockClient}

	message := entities.Message{
		Type:    "TestType",
		Payload: make(chan int),
	}

	err := snsInstance.Publish(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to serialize message for logging")
}
