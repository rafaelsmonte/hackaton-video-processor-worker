package adapters

type IVideoProcessorMessaging interface {
	Publish(message string) error
}
