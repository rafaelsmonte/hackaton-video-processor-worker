package adapters

import (
	"hackaton-video-processor-worker/internal/domain/entities"
)

type IVideoProcessorMessaging interface {
	Publish(message entities.Message) error
}
