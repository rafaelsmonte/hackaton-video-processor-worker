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

type SQSClient interface {
	ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

type SQSService struct {
	SqsClient SQSClient
	QueueURL  string
	Handler   *AppHandlers
}

type AppHandlers struct {
	VideoProcessorHandler *handlers.VideoHandler
}

func NewSQSService(region, queueURL string, handler *AppHandlers) *SQSService {
	var awsConfig *aws.Config

	awsConfig = &aws.Config{
		Region: aws.String(region),
	}
	log.Println("Using AWS SQS")

	sess := session.Must(session.NewSession(awsConfig))

	sqsClient := sqs.New(sess)

	return &SQSService{
		SqsClient: sqsClient,
		QueueURL:  queueURL,
		Handler:   handler,
	}
}

func (s *SQSService) StartConsuming(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down SQS consumer...")
			return
		default:
			msgs, err := s.ReceiveMessages()
			go func() {

				if err != nil {
					log.Fatalf("Error receiving messages: %v", err)
					return
				}

				for _, msg := range msgs {
					if err := s.ProcessMessage(msg); err != nil {
						log.Printf("Error processing message: %v", err)
						return
					}
				}
			}()

		}
	}
}

func (s *SQSService) ReceiveMessages() ([]*sqs.Message, error) {
	output, err := s.SqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.QueueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(1),
	})
	if err != nil {
		return nil, err
	}
	return output.Messages, nil
}

func (s *SQSService) ProcessMessage(msg *sqs.Message) error {
	log.Printf("Processing message: %s", aws.StringValue(msg.Body))

	if err := s.Handler.VideoProcessorHandler.HandleMessage(msg.Body); err != nil {
		return fmt.Errorf("failed to handle message: %w", err)
	}

	if err := s.DeleteMessage(msg.ReceiptHandle); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *SQSService) DeleteMessage(receiptHandle *string) error {
	_, err := s.SqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.QueueURL),
		ReceiptHandle: receiptHandle,
	})
	return err
}

func SetUpSQSService() {
	ctx := context.Background()
	queueURL := os.Getenv("SQS_QUEUE_URL")
	dlqURL := os.Getenv("SQS_DLQ_URL")
	region := os.Getenv("AWS_REGION")
	processingHandlers := ConfigProcessingHandlers()
	dlqHandlers := ConfigDLQHandlers()
	sqsService := NewSQSService(region, queueURL, processingHandlers)
	sqsServiceDLQ := NewSQSService(region, dlqURL, dlqHandlers)

	log.Println("Starting SQS consumer...")
	go sqsService.StartConsuming(ctx)
	sqsServiceDLQ.StartConsuming(ctx)
}
