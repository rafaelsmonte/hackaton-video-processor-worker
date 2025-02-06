//go:generate mockgen -destination=mocks/mock_video_processor_converter.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorConverter
//go:generate mockgen -destination=mocks/mock_video_processor_messaging.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorMessaging
//go:generate mockgen -destination=mocks/mock_video_processor_storage.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorStorage
//go:generate mockgen -destination=mocks/mock_video_processor_compressor.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorCompressor
//go:generate mockgen -destination=mocks/sqs_mock.go -package=sqs github.com/aws/aws-sdk-go-v2/service/sqs SQSAPI

package usecases
