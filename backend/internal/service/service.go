package service

import (
	"fmt"

	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
)

type ShellError struct {
	Msg      string
	ExitCode int
}

func (e *ShellError) Error() string {
	return fmt.Sprintf("ExitCode %d: %s", e.ExitCode, e.Msg)
}

type ApiServiceInterface interface {
	DockerPs(host string) (map[string][]string, error)
	CloneAndBuild(url, pjtName, contextPath string) (map[string][]string, error)
	RunContainer(host *request.PostRunProjectRequest) (map[string][]string, error)
	StopContainer(req *request.PostDockerStopRequest) (map[string][]string, error)
	RestartContainer(req *request.PostDockerRestartRequest) (map[string][]string, error)
	RemoveContainer(req *request.PostDockerRemoveRequest) (map[string][]string, error)
}
