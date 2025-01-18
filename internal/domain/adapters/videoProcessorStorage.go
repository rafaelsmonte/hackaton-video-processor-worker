package adapters

import "hackaton-video-processor-worker/internal/domain/entities"

type IVideoProcessorStorage interface {
	Download(entities.File) (entities.File, error)
	Upload(entities.File) (string, error)
}
