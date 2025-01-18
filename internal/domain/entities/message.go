package entities

const (
	StartProcessingMessage MessageTypeEnum = "MSG_EXTRACT_SNAPSHOT_STARTED"
	ExtractErrorMessage    MessageTypeEnum = "MSG_EXTRACT_SNAPSHOT_ERROR"
	ExtractSuccessMessage  MessageTypeEnum = "MSG_EXTRACT_SNAPSHOT_SUCCESS"
	SendErrorMessage       MessageTypeEnum = "MSG_SEND_SNAPSHOT_EXTRACTION_ERROR"
	SendSuccessMessage     MessageTypeEnum = "MSG_SEND_SNAPSHOT_EXTRACTION_SUCCESS"
	sender                 string          = "VIDEO_IMAGE_PROCESSOR_SERVICE"
	TargetVideoAPIService  Target          = "VIDEO_API_SERVICE"
	TargetEmailService     Target          = "EMAIL_SERVICE"
)

type MessageTypeEnum string
type Target string

type Message struct {
	Sender       string
	Target       Target
	MessatgeType MessageTypeEnum
	Payload      interface{}
}

type ExtractSuccessPayload struct {
	VideoSnapshotsUrl string
	VideoId           string
}
type ExtractErrorPayload struct {
	VideoId          string
	ErrorMessage     string
	ErrorDescription string
}

type ExtractSendSuccessPayload struct {
	UserID            string
	VideoSnapshotsUrl string
	VideoUrl          string
	VideoDescription  string
}
type ExtractSendErrorPayload struct {
	UserID            string
	VideoSnapshotsUrl string
	VideoUrl          string
	VideoDescription  string
	ErrorMessage      string
	ErrorDescription  string
}

func NewMessage(target Target, messageType MessageTypeEnum, payload interface{}) Message {
	return Message{
		Sender:       sender,
		Target:       target,
		MessatgeType: messageType,
		Payload:      payload,
	}
}
