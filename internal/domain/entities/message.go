package entities

const (
	StartProcessingMessage MessageTypeEnum = "MSG_EXTRACT_SNAPSHOT_PROCESSING"
	ExtractErrorMessage    MessageTypeEnum = "MSG_EXTRACT_SNAPSHOT_ERROR"
	ExtractSuccessMessage  MessageTypeEnum = "MSG_EXTRACT_SNAPSHOT_SUCCESS"
	sender                 string          = "VIDEO_IMAGE_PROCESSOR_SERVICE"
	TargetVideoSQSService  Target          = "VIDEO_API_SERVICE"
)

type MessageTypeEnum string
type Target string

type Message struct {
	Sender  string          `json:"sender"`
	Target  Target          `json:"target"`
	Type    MessageTypeEnum `json:"type"`
	Payload interface{}     `json:"payload"`
}

type ExtractSuccessPayload struct {
	VideoSnapshotsUrl string `json:"videoSnapshotsUrl"`
	VideoId           string `json:"videoId"`
	UserId            string `json:"userId"`
}
type ExtractErrorPayload struct {
	VideoId          string `json:"videoId"`
	UserId           string `json:"userId"`
	ErrorMessage     string `json:"errorMessage"`
	ErrorDescription string `json:"errorDescription"`
}

type StartProcessingPayload struct {
	VideoId string `json:"videoId"`
	UserId  string `json:"userId"`
}
type GenericErrorPayload struct {
	VideoId string `json:"videoId"`
	UserId  string `json:"userId"`
}

func NewMessage(target Target, messageType MessageTypeEnum, payload interface{}) Message {
	return Message{
		Sender:  sender,
		Target:  target,
		Type:    messageType,
		Payload: payload,
	}
}
