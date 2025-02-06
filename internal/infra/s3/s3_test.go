package s3

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"hackaton-video-processor-worker/internal/domain/entities"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func (m *MockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

type ReadCloserWrapper struct {
	*bytes.Reader
}

func (r *ReadCloserWrapper) Close() error {
	return nil
}

func NewReadCloserWrapper(b []byte) *ReadCloserWrapper {
	return &ReadCloserWrapper{Reader: bytes.NewReader(b)}
}
func TestUpload(t *testing.T) {
	mockS3Client := new(MockS3Client)
	s3Instance := &S3{
		Client: mockS3Client,
		Region: "us-east-1",
	}
	os.Setenv("S3_IMAGES_BUCKET_NAME", "test-bucket")
	file := entities.File{Id: "test.txt", Content: []byte{12}, Name: "test.txt", UserId: "123"}
	os.Create(file.Name)
	defer os.Remove(file.Name)

	mockS3Client.On("PutObject", mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == file.UserId+"/"+file.Name
	})).Return(&s3.PutObjectOutput{}, nil)

	uploadURL, err := s3Instance.Upload(file)

	assert.NoError(t, err)
	assert.Equal(t, uploadURL, "https://test-bucket.s3.us-east-1.amazonaws.com/123/test.txt")

	mockS3Client.AssertExpectations(t)
}

func TestUpload_Error(t *testing.T) {
	mockS3Client := new(MockS3Client)
	s3Instance := &S3{
		Client: mockS3Client,
		Region: "us-east-1",
	}
	os.Setenv("S3_IMAGES_BUCKET_NAME", "test-bucket")
	file := entities.File{Id: "test.txt", Content: []byte{12}, Name: "test.txt", UserId: "123"}

	os.Create(file.Id)
	defer os.Remove(file.Id)

	mockS3Client.On("PutObject", mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == file.UserId+"/"+file.Name
	})).Return(&s3.PutObjectOutput{}, errors.New("ERROR"))

	uploadURL, err := s3Instance.Upload(file)

	assert.Error(t, err)
	assert.Empty(t, uploadURL)

	mockS3Client.AssertExpectations(t)
}

func TestDownload(t *testing.T) {
	mockS3Client := new(MockS3Client)
	s3Instance := &S3{
		Client: mockS3Client,
		Region: "us-east-1",
	}
	os.Setenv("S3_VIDEO_BUCKET_NAME", "test-bucket")
	file := entities.File{Id: "test.txt", Content: []byte{12}, Name: "test.txt", UserId: "123"}

	os.Create(file.Id)
	defer os.Remove(file.Id)

	mockS3Client.On("GetObject", mock.Anything, mock.MatchedBy(func(input *s3.GetObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == file.UserId+"/"+file.Name
	})).Return(&s3.GetObjectOutput{
		Body: NewReadCloserWrapper([]byte("file content")),
	}, nil)

	downloadedFile, err := s3Instance.Download(file)

	assert.NoError(t, err)
	assert.Equal(t, []byte("file content"), downloadedFile.Content)

	mockS3Client.AssertExpectations(t)
}

func TestDownload_Error(t *testing.T) {
	mockS3Client := new(MockS3Client)
	s3Instance := &S3{
		Client: mockS3Client,
		Region: "us-east-1",
	}
	os.Setenv("S3_VIDEO_BUCKET_NAME", "test-bucket")
	file := entities.File{Id: "test.txt", Content: []byte{12}, Name: "test.txt", UserId: "123"}

	os.Create(file.Id)
	defer os.Remove(file.Id)

	mockS3Client.On("GetObject", mock.Anything, mock.MatchedBy(func(input *s3.GetObjectInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == file.UserId+"/"+file.Name
	})).Return(&s3.GetObjectOutput{
		Body: NewReadCloserWrapper([]byte("file content")),
	}, errors.New("unable to upload"))

	downloadedFile, err := s3Instance.Download(file)

	assert.Error(t, err)
	assert.Empty(t, downloadedFile)

	mockS3Client.AssertExpectations(t)

}
