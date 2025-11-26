package request

type PostDockerPsRequest struct {
	Host string `json:"host"`
}

type PostBuildProjectRequest struct {
	PjtName     string `json:"pjt_name"`
	URL         string `json:"url"`
	ContextPath string `json:"context_path"`
}

type PostRunProjectRequest struct {
	Host           string   `json:"host"`
	Image          string   `json:"image"`
	PortForwarding []string `json:"port"`
	Name           string   `json:"name"`
	Volume         []string `json:"volume"`
	Env            []string `json:"env"`
}
