package sns

import (
	"context"
	"encoding/json"
	"errors"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
)

// MockSNSClient implements SNSClient interface for testing
type MockSNSClient struct {
	publishFunc func(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func (m *MockSNSClient) Publish(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return m.publishFunc(ctx, input, optFns...)
}

func TestSNS_Publish(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		message     entities.Message
		setupMock   func() *MockSNSClient
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
	}{
		{
			name: "successful publish",
			message: entities.Message{
				Sender: "test-sender",
				Target: "VIDEO_API_SERVICE", // Preencha com o valor apropriado
				Type:   "test-type",         // Altere para o tipo correto
				Payload: map[string]interface{}{
					"status": "processing",
				},
			},
			setupMock: func() *MockSNSClient {
				return &MockSNSClient{
					publishFunc: func(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
						return &sns.PublishOutput{}, nil
					},
				}
			},
			setupEnv: func() {
				os.Setenv("SNS_TOPIC_ARN", "test-topic-arn")
			},
			cleanupEnv: func() {
				os.Unsetenv("SNS_TOPIC_ARN")
			},
			expectError: false,
		},
		{
			name: "missing topic ARN",
			message: entities.Message{
				Sender: "test-sender",
				Target: "VIDEO_API_SERVICE", // Preencha com o valor apropriado
				Type:   "test-type",         // Altere para o tipo correto
				Payload: map[string]interface{}{
					"status": "processing",
				},
			},
			setupMock: func() *MockSNSClient {
				return &MockSNSClient{
					publishFunc: func(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
						return &sns.PublishOutput{}, nil
					},
				}
			},
			setupEnv: func() {
				os.Unsetenv("SNS_TOPIC_ARN")
			},
			cleanupEnv:  func() {},
			expectError: true,
		},
		{
			name: "publish error",
			message: entities.Message{
				Sender: "test-sender",
				Target: "VIDEO_API_SERVICE", // Preencha com o valor apropriado
				Type:   "test-type",         // Altere para o tipo correto
				Payload: map[string]interface{}{
					"status": "processing",
				},
			},
			setupMock: func() *MockSNSClient {
				return &MockSNSClient{
					publishFunc: func(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
						return nil, errors.New("publish error")
					},
				}
			},
			setupEnv: func() {
				os.Setenv("SNS_TOPIC_ARN", "test-topic-arn")
			},
			cleanupEnv: func() {
				os.Unsetenv("SNS_TOPIC_ARN")
			},
			expectError: true,
		},
		{
			name: "json marshal error",
			message: entities.Message{
				Sender: "test-sender",
				Target: "VIDEO_API_SERVICE", // Preencha com o valor apropriado
				Type:   "test-type",         // Altere para o tipo correto
				Payload: map[string]interface{}{
					"channel": make(chan int), // Isso causará falha na serialização JSON
				},
			},
			setupMock: func() *MockSNSClient {
				return &MockSNSClient{
					publishFunc: func(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
						return &sns.PublishOutput{}, nil
					},
				}
			},
			setupEnv: func() {
				os.Setenv("SNS_TOPIC_ARN", "test-topic-arn")
			},
			cleanupEnv: func() {
				os.Unsetenv("SNS_TOPIC_ARN")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setupEnv()
			defer tt.cleanupEnv()

			mockClient := tt.setupMock()
			snsInstance := &SNS{
				Client: mockClient,
			}

			// Execute
			err := snsInstance.Publish(tt.message)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewSNS(t *testing.T) {
	t.Run("successful initialization", func(t *testing.T) {
		// This test might fail in environments without AWS credentials
		// You might want to mock the config.LoadDefaultConfig call in a production environment
		snsInstance, err := NewSNS()
		if err != nil {
			t.Skip("Skipping test due to AWS credentials not being available")
		}
		assert.NotNil(t, snsInstance)
		assert.NoError(t, err)
	})
}

// TestMessageJSONSerialization tests the JSON serialization of the Message struct
func TestMessageJSONSerialization(t *testing.T) {
	message := entities.Message{
		Sender: "test-sender",
		Target: "VIDEO_API_SERVICE", // Preencha com o valor apropriado
		Type:   "test-type",         // Altere para o tipo correto
		Payload: map[string]interface{}{
			"key": "value",
		},
	}

	// Test MarshalIndent
	_, err := json.MarshalIndent(message, "", "  ")
	assert.NoError(t, err)

	// Test Marshal
	_, err = json.Marshal(message)
	assert.NoError(t, err)
}
