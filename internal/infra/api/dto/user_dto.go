package dto

type VideoProcessRequest struct {
	Name string `query:"name"`
	Id   string `query:"id"`
	Path string `query:"path"`
}

type VideoProcessResponse struct {
	Message string `json:"message"`
}
