//go:generate mockgen -destination=mocks/mock_video_processor_converter.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorConverter
//go:generate mockgen -destination=mocks/mock_video_processor_messaging.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorMessaging
//go:generate mockgen -destination=mocks/mock_video_processor_storage.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorStorage
//go:generate mockgen -destination=mocks/mock_video_processor_compressor.go -package=mocks hackaton-video-processor-worker/internal/domain/adapters IVideoProcessorCompressor

package usecases
