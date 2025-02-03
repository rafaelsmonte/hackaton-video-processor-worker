package sqs

import (
	"context"
	"fmt"
	"hackaton-video-processor-worker/internal/infra/sqs/handlers"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSService struct {
	sqsClient *sqs.SQS
	queueURL  string
	handler   *AppHandlers
}

type AppHandlers struct {
	videoProcessorHandler *handlers.VideoHandler
}

func NewSQSService(region, queueURL string, handler *AppHandlers) *SQSService {
	env := os.Getenv("ENV")
	var awsConfig *aws.Config

	if env == "DEV" {
		awsConfig = &aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("http://localhost:4566"),
		}
		log.Println("Using LocalStack for SQS")
	} else {
		awsConfig = &aws.Config{
			Region: aws.String(region),
		}
		log.Println("Using AWS SQS")
	}

	sess := session.Must(session.NewSession(awsConfig))

	sqsClient := sqs.New(sess)

	return &SQSService{
		sqsClient: sqsClient,
		queueURL:  queueURL,
		handler:   handler,
	}
}

func (s *SQSService) StartConsuming(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down SQS consumer...")
			return
		default:
			msgs, err := s.receiveMessages()
			go func() {

				if err != nil {
					log.Fatalf("Error receiving messages: %v", err)
					return
				}

				for _, msg := range msgs {
					if err := s.processMessage(msg); err != nil {
						log.Printf("Error processing message: %v", err)
						return
					}
				}
			}()

		}
	}
}

func (s *SQSService) receiveMessages() ([]*sqs.Message, error) {
	output, err := s.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(5),
	})
	if err != nil {
		return nil, err
	}
	return output.Messages, nil
}

func (s *SQSService) processMessage(msg *sqs.Message) error {
	log.Printf("Processing message: %s", aws.StringValue(msg.Body))

	if err := s.handler.videoProcessorHandler.HandleMessage(msg.Body); err != nil {
		return fmt.Errorf("failed to handle message: %w", err)
	}

	if err := s.deleteMessage(msg.ReceiptHandle); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *SQSService) deleteMessage(receiptHandle *string) error {
	_, err := s.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: receiptHandle,
	})
	return err
}

func SetUpSQSService() {
	ctx := context.Background()
	queueURL := os.Getenv("SQS_QUEUE_URL")
	region := os.Getenv("AWS_REGION")
	handler := ConfigHandlers()
	sqsService := NewSQSService(region, queueURL, handler)

	log.Println("Starting SQS consumer...")
	sqsService.StartConsuming(ctx)
}
