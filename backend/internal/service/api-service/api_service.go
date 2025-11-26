package apiservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jaeun-Choi98/container-bay/internal/config"
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/Jaeun-Choi98/modules/shell"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

type ApiService struct {
	Cfg *config.Config
}

func NewApiService(cfg *config.Config) *ApiService {
	Init(cfg)
	return &ApiService{
		Cfg: cfg,
	}
}

func Init(cfg *config.Config) error {
	if err := DockerRepoLogin(cfg); err != nil {
		return fmt.Errorf("[API Service] failed to init")
	}

	return nil
}

func DockerRepoLogin(cfg *config.Config) error {
	executor := shell.NewScriptExecutor(cfg.ShellDir)
	script := fmt.Sprintf(`
		echo %s | sudo -S docker login %s:5000 -u %s -p %s
	`, cfg.BuildSvrPasswd, cfg.DockerRepoIp, cfg.DockerRepoId, cfg.DockerRepoPwd)
	result, err := executor.Execute(context.Background(), cfg.ShellName, "docker_login.sh", script, true)
	if err != nil {
		return fmt.Errorf("[API Service] failed to execute script(docker_login.sh): %v", err)
	}
	if result.Error != nil {
		return fmt.Errorf("[API Service] failed to execute script(docker_login.sh): %v", err)
	}
	if len(result.Stderr) != 0 {
		logger.Println("=== STDERR ===")
		for _, line := range result.Stderr {
			logger.Println(line)
		}
	}
	return nil
}

func (s *ApiService) DockerPs(host string) ([]string, error) {

	hostLocalSepSlash := filepath.FromSlash(host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDaemonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	log.Printf("[debug] docker daemon host: %s", dockerDaemonHost)

	executor := shell.NewScriptExecutor(s.Cfg.ShellDir)

	script := fmt.Sprintf(`
		echo %s | sudo -S docker -H %s ps -a
	`, s.Cfg.BuildSvrPasswd, dockerDaemonHost)

	result, err := executor.Execute(context.Background(), s.Cfg.ShellName, "docker_ps.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_ps.sh): %v", err), response.FAIL)
	}

	logger.Println("<=== Run docker_ps.sh ===>")
	logger.Printf("content: %s", script)

	var res []string

	logger.Println("=== Execution Result ===")
	res = append(res, "=== Execution Result ===")

	logger.Printf("Duration: %v\n", result.Duration)
	res = append(res, fmt.Sprintf("Duration: %v", result.Duration))

	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res = append(res, fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res = append(res, "=== STDOUT ===")

	for _, line := range result.Stdout {
		logger.Println(line)
		res = append(res, line)
	}
	logger.Println("=== STDERR ===")
	res = append(res, "=== STDERR ===")
	for _, line := range result.Stderr {
		logger.Println(line)
		res = append(res, line)
	}

	if result.Error != nil {
		logger.Printf("Error: %v\n", result.Error)
		res = append(res, fmt.Sprintf("Error: %v", result.Error))
	}
	logger.Println("<=== END ===>")

	return res, nil
}

func (s *ApiService) CloneAndBuild(url, pjtName, contextPath string) ([]string, error) {
	pjtPath := filepath.Join(s.Cfg.RepoDir, pjtName)
	//log.Printf("[debug] pjt_path: %s", pjtPath)

	// clone하기 전에, 해당 path에 폴더가 이미 있는지 확인.
	os.RemoveAll(pjtPath)

	// repo는 필요하지 않음.
	_, err := git.PlainClone(pjtPath, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: s.Cfg.GitId,
			Password: s.Cfg.GitPwd,
		},
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to clone: %v", err), response.FAIL)
	}

	buildPath := filepath.Join(pjtPath, contextPath)
	// log.Printf("[debug] build_path: %s", buildPath)

	executor := shell.NewScriptExecutor(s.Cfg.ShellDir)
	dockerRepoUrl := fmt.Sprintf("%s:%s", s.Cfg.DockerRepoIp, s.Cfg.DockerRepoPort)
	script := fmt.Sprintf(`
		echo %s | sudo -S docker build -t %s:latest %s
		echo %s | sudo -S docker tag %s:latest %s/%s:latest
		echo %s | sudo -S docker push %s/%s:latest
		echo %s | sudo -S docker system prune -f
	`, s.Cfg.BuildSvrPasswd, pjtName, buildPath,
		s.Cfg.BuildSvrPasswd, pjtName, dockerRepoUrl, pjtName,
		s.Cfg.BuildSvrPasswd, dockerRepoUrl, pjtName,
		s.Cfg.BuildSvrPasswd)

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	rst, err := executor.Execute(context.Background(), s.Cfg.ShellName, "docker_build.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_build.sh): %v", err), response.FAIL)
	}

	logger.Println("<=== Run docker_build.sh ===>")
	logger.Printf("content: %s", script)

	var res []string

	logger.Println("=== Execution Result ===")
	res = append(res, "=== Execution Result ===")

	logger.Printf("Duration: %v\n", rst.Duration)
	res = append(res, fmt.Sprintf("Duration: %v", rst.Duration))

	logger.Printf("Exit Code: %d\n", rst.ExitCode)
	res = append(res, fmt.Sprintf("Exit Code: %d", rst.ExitCode))

	logger.Println("=== STDOUT ===")
	res = append(res, "=== STDOUT ===")

	for _, line := range rst.Stdout {
		logger.Println(line)
		res = append(res, line)
	}
	logger.Println("=== STDERR ===")
	res = append(res, "=== STDERR ===")
	for _, line := range rst.Stderr {
		logger.Println(line)
		res = append(res, line)
	}

	if rst.Error != nil {
		logger.Printf("Error: %v\n", rst.Error)
		res = append(res, fmt.Sprintf("Error: %v", rst.Error))
	}
	logger.Println("<=== END ===>")
	// 빌드가 끝났다면 pjt 삭제
	os.RemoveAll(pjtPath)

	return res, nil
}

func (s *ApiService) RunContainer(req *request.PostRunProjectRequest) ([]string, error) {
	hostLocalSepSlash := filepath.FromSlash(req.Host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDaemonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	//log.Printf("[debug] docker daemon host: %s", dockerDaemonHost)

	execurator := shell.NewScriptExecutor(s.Cfg.ShellDir)
	var script strings.Builder
	script.WriteString(fmt.Sprintf(`
			echo %s | sudo -S docker -H %s run -d `, s.Cfg.BuildSvrPasswd, dockerDaemonHost,
	))

	for _, p := range req.PortForwarding {
		script.WriteString("-p ")
		script.WriteString(fmt.Sprintf("%s ", p))
	}

	script.WriteString("-v ")
	for _, v := range req.Volume {
		script.WriteString(fmt.Sprintf("%s ", v))
	}

	script.WriteString("-e ")
	for _, e := range req.Env {
		script.WriteString(fmt.Sprintf("%s ", e))
	}

	script.WriteString(fmt.Sprintf("--name %s %s:%s/%s:latest", req.Name, s.Cfg.DockerRepoIp, s.Cfg.DockerRepoPort, req.Image))

	rst, err := execurator.Execute(context.Background(), s.Cfg.ShellName, "docker_run.sh", script.String(), true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_run.sh): %v", err), response.FAIL)
	}

	logger.Println("<=== Run docker_run.sh ===>")
	logger.Printf("content: %s", script.String())

	var res []string

	logger.Println("=== Execution Result ===")
	res = append(res, "=== Execution Result ===")

	logger.Printf("Duration: %v\n", rst.Duration)
	res = append(res, fmt.Sprintf("Duration: %v", rst.Duration))

	logger.Printf("Exit Code: %d\n", rst.ExitCode)
	res = append(res, fmt.Sprintf("Exit Code: %d", rst.ExitCode))

	logger.Println("=== STDOUT ===")
	res = append(res, "=== STDOUT ===")

	for _, line := range rst.Stdout {
		logger.Println(line)
		res = append(res, line)
	}
	logger.Println("=== STDERR ===")
	res = append(res, "=== STDERR ===")
	for _, line := range rst.Stderr {
		logger.Println(line)
		res = append(res, line)
	}

	if rst.Error != nil {
		logger.Printf("Error: %v\n", rst.Error)
		res = append(res, fmt.Sprintf("Error: %v", rst.Error))
	}

	logger.Println("<=== END ===>")
	return res, nil
}
