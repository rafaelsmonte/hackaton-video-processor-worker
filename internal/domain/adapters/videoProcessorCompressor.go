package adapters

import "hackaton-video-processor-worker/internal/domain/entities"

type IVideoProcessorCompressor interface {
	Compress(entities.Folder) (entities.File, error)
}
