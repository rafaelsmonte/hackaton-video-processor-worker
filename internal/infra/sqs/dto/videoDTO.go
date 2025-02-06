package dto

type VideoProcessRequest struct {
	VideoId   string `json:"videoId"`
	UserId    string `json:"userId"`
	VideoName string `json:"videoName"`
}

type VideoProcessResponse struct {
	Message string `json:"message"`
}
