package request

type PostImageLsRequest struct {
	Host string `json:"host"`
}

type PostImageRmRequest struct {
	Host      string `json:"host"`
	ImageName string `json:"name"`
}
