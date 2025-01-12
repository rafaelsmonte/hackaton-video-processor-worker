package handlers

import (
	"encoding/json"
	"hackaton-video-processor-worker/internal/domain/usecases"
	"hackaton-video-processor-worker/internal/infra/api/dto"
	"hackaton-video-processor-worker/internal/infra/api/mappers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	videoUsecase usecases.IConvertVideoUsecase
}

func NewVideoHandler(videoUsecase usecases.IConvertVideoUsecase) *VideoHandler {
	return &VideoHandler{
		videoUsecase: videoUsecase,
	}
}

func (h *VideoHandler) Handle(c echo.Context) error {
	var video dto.VideoProcessRequest

	if err := c.Bind(&video); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	input := mappers.ProcessVideoInput(video)
	output, err := h.videoUsecase.Execute(input)
	if err != nil {
		c.Error(err)
	}

	response := mappers.ProcessVideoResponse(output)
	responseJson, err := json.Marshal(response)
	if err != nil {
		c.Error(err)
	}

	return c.String(http.StatusOK, string(responseJson))
}
