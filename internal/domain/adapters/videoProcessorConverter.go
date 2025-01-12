package adapters

import "hackaton-video-processor-worker/internal/domain/entities"

type IVideoProcessorConverter interface {
	ConvertToImages(entities.File) (entities.Folder, error)
}
