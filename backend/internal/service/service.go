package service

import "github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"

type ApiServiceInterface interface {
	DockerPs(host string) ([]string, error)
	CloneAndBuild(url, pjtName, contextPath string) ([]string, error)
	RunContainer(host *request.PostRunProjectRequest) ([]string, error)
}
