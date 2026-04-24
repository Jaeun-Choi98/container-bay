package request

type PostAddDaemonRequest struct {
	Host  string `json:"host"`
	Label string `json:"label"`
}
