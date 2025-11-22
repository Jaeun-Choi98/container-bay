package request

type PostBuildProjectRequest struct {
	PjtName     string `json:"pjt_name"`
	URL         string `json:"url"`
	ContextPath string `json:"context_path"`
}
