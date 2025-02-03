package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSClient interface {
	Publish(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
type SNS struct {
	Client SNSClient
}

func (snsInstance *SNS) Publish(message entities.Message) error {

	log.Println(string(message.MessatgeType), message.Payload)

	ctx := context.Background()
	topicArn := os.Getenv("SNS_TOPIC_ARN")
	if topicArn == "" {
		return fmt.Errorf("SNS_TOPIC_ARN is not set")
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(messageBody)),
	}

	_, err = snsInstance.Client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message published successfully to topic %s\n", topicArn)
	return nil
}

func NewSNS() (adapters.IVideoProcessorMessaging, error) {
	var cfg aws.Config
	var err error

	if os.Getenv("ENV") == "DEV" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localhost:4566",
					SigningRegion: "us-east-1",
				}, nil
			})),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				"test",
				"test",
				"",
			)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS SDK config: %w", err)
		}
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS SDK config: %w", err)
		}
	}

	return &SNS{
		Client: sns.NewFromConfig(cfg),
	}, nil
}
