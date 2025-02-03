package s3

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Client *s3.Client
	Region string
}

func NewS3() (adapters.IVideoProcessorStorage, error) {
	env := os.Getenv("ENV")
	region := os.Getenv("S3_REGION")
	var cfg aws.Config
	var err error

	if env == "DEV" {
		// Use LocalStack when in DEV environment
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					// Set LocalStack endpoint for S3
					return aws.Endpoint{
						URL:               "http://localhost:4566", // LocalStack default S3 endpoint
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown service %s", service)
			})),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3{
		Client: client,
		Region: region,
	}, nil
}

func (s3Instance *S3) Download(file entities.File) (entities.File, error) {
	bucketName := os.Getenv("S3_VIDEO_BUCKET_NAME")
	fmt.Println(bucketName, file.Id)

	getObjectRequest := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Id),
	}

	resp, err := s3Instance.Client.GetObject(context.TODO(), getObjectRequest)
	if err != nil {
		return entities.File{}, fmt.Errorf("unable to download file, %v", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return entities.File{}, fmt.Errorf("unable to read file content, %v", err)
	}

	file.Content = buf.Bytes()

	return file, nil
}

func (s3Instance *S3) Upload(file entities.File) (string, error) {
	bucketName := os.Getenv("S3_IMAGES_BUCKET_NAME")

	fileContent, err := os.Open(file.Id)
	if err != nil {
		return "", fmt.Errorf("unable to open file, %v", err)
	}
	defer fileContent.Close()

	putObjectRequest := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Id),
		Body:   fileContent,
	}

	_, err = s3Instance.Client.PutObject(context.TODO(), putObjectRequest)
	if err != nil {
		return "", fmt.Errorf("unable to upload file, %v", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, s3Instance.Region, file.Id), nil
}
