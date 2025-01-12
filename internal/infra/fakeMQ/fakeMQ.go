package fakeMQ

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"log"
)

type FakeMQ struct {
}

func (f *FakeMQ) Publish(message string) error {
	log.Println("Mensagem publicada na fakeMQ - ", message)
	return nil
}

func NewFakeMQ() adapters.IVideoProcessorMessaging {
	return &FakeMQ{}
}
