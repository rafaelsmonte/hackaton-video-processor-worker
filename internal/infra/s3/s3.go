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
}

func (f *S3) Download(file entities.File) (entities.File, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return entities.File{}, fmt.Errorf("unable to load SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	bucketName := "your-bucket-name"
	key := file.Name

	output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
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
func (f *S3) Upload(file entities.File) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	filePath := filepath.Join(file.Path, file.Name)
	bucketName := "your-bucket-name"
	key := "example.txt"

	rFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file, %v", err)
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(rFile),
	})
	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
	}

	fmt.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, key)
	return nil
}

func NewS3() adapters.IVideoProcessorStorage {
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	log.Fatalf("unable to load SDK config, %v", err)
	// }

	// // Create an S3 client
	// s3Client := s3.NewFromConfig(cfg)

	// // List S3 buckets
	// listBucketsOutput, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	// if err != nil {
	// 	log.Fatalf("unable to list buckets, %v", err)
	// }

	// fmt.Println("Buckets:")
	// for _, bucket := range listBucketsOutput.Buckets {
	// 	fmt.Printf("* %s (created on %v)\n", aws.ToString(bucket.Name), bucket.CreationDate)
	// }
	return &S3{}
}
