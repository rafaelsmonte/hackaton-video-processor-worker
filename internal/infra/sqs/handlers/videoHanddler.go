package handlers

import (
	"encoding/json"
	"fmt"
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/sqs/dto"
	"hackaton-video-processor-worker/internal/infra/sqs/mappers"
	"log"
)

type VideoHandler struct {
	videoUsecase usecases.IConvertVideoUsecase
}

func NewVideoHandler(videoUsecase usecases.IConvertVideoUsecase) *VideoHandler {
	return &VideoHandler{
		videoUsecase: videoUsecase,
	}
}

func (h *VideoHandler) HandleMessage(body *string) error {
	var message struct {
		Sender  string `json:"sender"`
		Target  string `json:"target"`
		Type    string `json:"type"`
		Payload struct {
			VideoUrl string `json:"videoUrl"`
			VideoId  string `json:"videoId"`
		} `json:"payload"`
	}
	fmt.Println("123")

	if err := json.Unmarshal([]byte(*body), &message); err != nil {
		log.Printf("Failed to unmarshal message body: %v", err)
		return fmt.Errorf("invalid message format: %w", err)
	}

	if message.Type != "MSG_EXTRACT_SNAPSHOT" {
		log.Printf("Unsupported message type: %s", message.Type)
		return fmt.Errorf("unsupported message type: %s", message.Type)
	}

	videoRequest := dto.VideoProcessRequest{
		VideoUrl: message.Payload.VideoUrl,
		VideoId:  message.Payload.VideoId,
	}

	input := mappers.ProcessVideoInput(videoRequest)
	output, err := h.videoUsecase.Execute(input)
	if err != nil {
		log.Printf("Failed to process video: %v", err)
		return fmt.Errorf("video processing error: %w", err)
	}

	response := mappers.ProcessVideoResponse(output)
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return fmt.Errorf("response marshalling error: %w", err)
	}

	log.Printf("Successfully processed video. Response: %s", string(responseJson))
	return nil
}
