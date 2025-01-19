package dto

type VideoProcessRequest struct {
	VideoUrl string `json:"videoUrl"`
	VideoId  string `json:"videoId"`
}

type VideoProcessResponse struct {
	Message string `json:"message"`
}
