package entities

type Message struct {
	Sender       string
	Target       string
	MessatgeType string
	Payload      payload
}

type payload struct {
	VideoUrl string
	VideoId  string
}
