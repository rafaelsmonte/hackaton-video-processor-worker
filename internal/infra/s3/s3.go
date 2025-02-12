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

type S3ClientInterface interface {
	GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}
type S3 struct {
	Client S3ClientInterface
	Region string
}

func NewS3() (adapters.IVideoProcessorStorage, error) {
	region := os.Getenv("S3_REGION")
	var cfg aws.Config
	var err error

	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

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
	fmt.Println(bucketName, file.Name)

	getObjectRequest := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.UserId + "/" + file.Name),
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

	fileContent, err := os.Open(file.Name)
	if err != nil {
		return "", fmt.Errorf("unable to open file, %v", err)
	}
	defer fileContent.Close()

	putObjectRequest := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.UserId + "/" + file.Name),
		Body:   fileContent,
	}

	_, err = s3Instance.Client.PutObject(context.TODO(), putObjectRequest)
	if err != nil {
		return "", fmt.Errorf("unable to upload file, %v", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, s3Instance.Region, file.UserId+"/"+file.Name), nil
}
