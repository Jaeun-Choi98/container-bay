package apiservice

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/service"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/Jaeun-Choi98/modules/shell"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

func (s *ApiService) DockerPs(host string) (map[string][]string, error) {

	hostLocalSepSlash := filepath.FromSlash(host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDaemonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	// log.Printf("[debug] docker daemon host: %s", dockerDaemonHost)

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

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execute_result"] = append(res["execute_result"], "=== Execution Result ===")

	logger.Printf("Duration: %v\n", result.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", result.Duration))

	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")

	// ps의 경우, Port 컬럼은 없을 경우 생략됨 -> port(5번째 인덱스)
	for lineIndex, line := range result.Stdout {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		idx := 0
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
				idx++
			}
			if lineIndex != 0 && trimStr == "" && idx == 5 {
				parseStr.WriteString(" ")
				parseStr.WriteString(";")
				idx++
			}
		}
		res["stdout"] = append(res["stdout"], parseStr.String())
	}

	logger.Println("=== STDERR ===")
	res["stderr"] = append(res["stderr"], "=== STDERR ===")
	for _, line := range result.Stderr {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stderr"] = append(res["stderr"], parseStr.String())
	}

	if result.Error != nil {
		logger.Printf("Error: %v\n", result.Error)
		res["stderr"] = append(res["stderr"], fmt.Sprintf("Error: %v", result.Error))
	}
	logger.Println("<=== END ===>")

	// 정상적으로 처리하지 못함
	if result.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: result.ExitCode,
			Msg:      "occured issue in processing docker_ps.sh",
		}
	}

	return res, nil
}

func (s *ApiService) CloneAndBuild(url, pjtName, contextPath string) (map[string][]string, error) {
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

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execute_result"] = append(res["execute_result"], "=== Execution Result ===")

	logger.Printf("Duration: %v\n", rst.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", rst.Duration))

	logger.Printf("Exit Code: %d\n", rst.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", rst.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")

	for _, line := range rst.Stdout {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stdout"] = append(res["stdout"], parseStr.String())
	}
	logger.Println("=== STDERR ===")
	res["stderr"] = append(res["stderr"], "=== STDERR ===")
	for _, line := range rst.Stderr {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stderr"] = append(res["stderr"], parseStr.String())
	}

	if rst.Error != nil {
		logger.Printf("Error: %v\n", rst.Error)
		res["stderr"] = append(res["stderr"], fmt.Sprintf("Error: %v", rst.Error))
	}
	logger.Println("<=== END ===>")
	// 빌드가 끝났다면 pjt 삭제
	os.RemoveAll(pjtPath)

	// 정상적으로 처리하지 못함
	if rst.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: rst.ExitCode,
			Msg:      "occured issue in processing script(docker_build.sh)",
		}
	}

	return res, nil
}

func (s *ApiService) RunContainer(req *request.PostRunProjectRequest) (map[string][]string, error) {
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

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execute_result"] = append(res["execute_result"], "=== Execution Result ===")

	logger.Printf("Duration: %v\n", rst.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", rst.Duration))

	logger.Printf("Exit Code: %d\n", rst.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", rst.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")

	for _, line := range rst.Stdout {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stdout"] = append(res["stdout"], parseStr.String())
	}
	logger.Println("=== STDERR ===")
	res["stderr"] = append(res["stderr"], "=== STDERR ===")
	for _, line := range rst.Stderr {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stderr"] = append(res["stderr"], parseStr.String())
	}

	if rst.Error != nil {
		logger.Printf("Error: %v\n", rst.Error)
		res["stderr"] = append(res["stderr"], fmt.Sprintf("Error: %v", rst.Error))
	}

	logger.Println("<=== END ===>")

	// 정상적으로 처리하지 못함
	if rst.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: rst.ExitCode,
			Msg:      "occured issue in processing script(docker_run.sh)",
		}
	}
	return res, nil
}

func (a *ApiService) StopContainer(req *request.PostDockerStopRequest) (map[string][]string, error) {
	hostLocalSepSlash := filepath.FromSlash(req.Host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDaemonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	executor := shell.NewScriptExecutor(a.Cfg.ShellDir)

	script := fmt.Sprintf(`
		echo %s | sudo -S docker -H %s stop %s
	`, a.Cfg.BuildSvrPasswd, dockerDaemonHost, req.ContainerName)

	result, err := executor.Execute(context.Background(), a.Cfg.ShellName, "docker_stop.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_stop.sh): %v", err),
			response.FAIL)
	}

	logger.Println("<=== Start docker_stop.sh ===>")
	logger.Printf("content: %s", script)

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execute_result"] = append(res["execute_result"], "=== Execution Result ===")

	logger.Printf("Duration: %v\n", result.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", result.Duration))

	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")

	// ps의 경우, Port 컬럼은 없을 경우 생략됨 -> port(5번째 인덱스)
	for _, line := range result.Stdout {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stdout"] = append(res["stdout"], parseStr.String())
	}

	logger.Println("=== STDERR ===")
	res["stderr"] = append(res["stderr"], "=== STDERR ===")
	for _, line := range result.Stderr {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stderr"] = append(res["stderr"], parseStr.String())
	}

	if result.Error != nil {
		logger.Printf("Error: %v\n", result.Error)
		res["stderr"] = append(res["stderr"], fmt.Sprintf("Error: %v", result.Error))
	}
	logger.Println("<=== END ===>")

	// 정상적으로 처리하지 못함
	if result.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: result.ExitCode,
			Msg:      "occured issue in processing docker_stop.sh",
		}
	}

	return res, nil
}

func (a *ApiService) RestartContainer(req *request.PostDockerRestartRequest) (map[string][]string, error) {
	hostLocalSepSlash := filepath.FromSlash(req.Host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDaemonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	executor := shell.NewScriptExecutor(a.Cfg.ShellDir)

	script := fmt.Sprintf(`
		echo %s | sudo -S docker -H %s rm %s
	`, a.Cfg.BuildSvrPasswd, dockerDaemonHost, req.ContainerName)

	result, err := executor.Execute(context.Background(), a.Cfg.ShellName, "docker_restart.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_restart.sh): %v", err),
			response.FAIL)
	}

	logger.Println("<=== Start docker_restart.sh ===>")
	logger.Printf("content: %s", script)

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execute_result"] = append(res["execute_result"], "=== Execution Result ===")

	logger.Printf("Duration: %v\n", result.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", result.Duration))

	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")

	// ps의 경우, Port 컬럼은 없을 경우 생략됨 -> port(5번째 인덱스)
	for _, line := range result.Stdout {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stdout"] = append(res["stdout"], parseStr.String())
	}

	logger.Println("=== STDERR ===")
	res["stderr"] = append(res["stderr"], "=== STDERR ===")
	for _, line := range result.Stderr {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stderr"] = append(res["stderr"], parseStr.String())
	}

	if result.Error != nil {
		logger.Printf("Error: %v\n", result.Error)
		res["stderr"] = append(res["stderr"], fmt.Sprintf("Error: %v", result.Error))
	}
	logger.Println("<=== END ===>")

	// 정상적으로 처리하지 못함
	if result.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: result.ExitCode,
			Msg:      "occured issue in processing docker_rm.sh",
		}
	}

	return res, nil
}

func (a *ApiService) RemoveContainer(req *request.PostDockerRemoveRequest) (map[string][]string, error) {
	hostLocalSepSlash := filepath.FromSlash(req.Host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDaemonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	executor := shell.NewScriptExecutor(a.Cfg.ShellDir)

	script := fmt.Sprintf(`
		echo %s | sudo -S docker -H %s rm %s
	`, a.Cfg.BuildSvrPasswd, dockerDaemonHost, req.ContainerName)

	result, err := executor.Execute(context.Background(), a.Cfg.ShellName, "docker_rm.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_rm.sh): %v", err),
			response.FAIL)
	}

	logger.Println("<=== Start docker_rm.sh ===>")
	logger.Printf("content: %s", script)

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execute_result"] = append(res["execute_result"], "=== Execution Result ===")

	logger.Printf("Duration: %v\n", result.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", result.Duration))

	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")

	// ps의 경우, Port 컬럼은 없을 경우 생략됨 -> port(5번째 인덱스)
	for _, line := range result.Stdout {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stdout"] = append(res["stdout"], parseStr.String())
	}

	logger.Println("=== STDERR ===")
	res["stderr"] = append(res["stderr"], "=== STDERR ===")
	for _, line := range result.Stderr {
		logger.Println(line)
		var parseStr strings.Builder
		splitStr := strings.Split(line, "  ")
		for _, str := range splitStr {
			trimStr := strings.TrimSpace(str)
			if trimStr != "" {
				parseStr.WriteString(trimStr)
				parseStr.WriteString(";")
			}
		}
		res["stderr"] = append(res["stderr"], parseStr.String())
	}

	if result.Error != nil {
		logger.Printf("Error: %v\n", result.Error)
		res["stderr"] = append(res["stderr"], fmt.Sprintf("Error: %v", result.Error))
	}
	logger.Println("<=== END ===>")

	// 정상적으로 처리하지 못함
	if result.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: result.ExitCode,
			Msg:      "occured issue in processing docker_rm.sh",
		}
	}

	return res, nil
}
