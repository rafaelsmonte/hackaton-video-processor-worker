package s3

import (
	"bytes"
	"context"
	"fmt"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Client *s3.Client
	Region string
}

func (s3Instance *S3) Download(file entities.File) (entities.File, error) {

	bucketName := "your-bucket-name"
	key := file.Name

	output, err := s3Instance.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer output.Body.Close()

	localFilePath := filepath.Join(file.Path, file.Name)

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	_, err = bytes.NewBuffer(nil).ReadFrom(output.Body)
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to write to local file: %w", err)
	}

	fmt.Printf("Successfully downloaded %s from bucket %s to %s\n", key, bucketName, localFilePath)

	return entities.NewFile(file.Id, file.Name, localFilePath), nil
}

// https://rafael-fiap.s3.us-east-1.amazonaws.com/the_number_of_the_beast.mp4
// export AWS_ACCESS_KEY_ID="your-access-key-id"
// export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
// export AWS_REGION="your-region"
func (s3Instance *S3) Upload(file entities.File) (string, error) {

	filePath := filepath.Join(file.Path, file.Name)
	bucketName := "your-bucket-name"
	key := "example.txt"

	rFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file, %v", err)
		return "", err

	}

	_, err = s3Instance.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(rFile),
	})
	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
		return "", err

	}

	fmt.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, key)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, s3Instance.Region, key)
	return url, nil
}

func NewS3() (adapters.IVideoProcessorStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return &S3{
		Client: s3Client,
		Region: cfg.Region,
	}, nil
}
