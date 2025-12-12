package apiservice

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/service"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/Jaeun-Choi98/modules/shell"
)

func (s *ApiService) ImageLs(host string) (map[string][]string, error) {
	hostLocalSepSlash := filepath.FromSlash(host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDeamonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	executor := shell.NewScriptExecutor(s.Cfg.ShellDir)
	script := fmt.Sprintf(`
		echo %s | sudo -S docker -H %s image ls
	`, s.Cfg.BuildSvrPasswd, dockerDeamonHost)

	result, err := executor.Execute(context.Background(), s.Cfg.ShellName, "docker_img_ls.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_img_ls.sh): %v", err), response.FAIL)
	}

	logger.Println("<=== Run docker_img_ls.sh ===>")
	logger.Printf("content: %s", script)

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execut_result"] = append(res["execut_result"], "=== Execution Result ===")
	logger.Printf("Duration: %v\n", result.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", result.Duration))
	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")
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

	if result.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: result.ExitCode,
			Msg:      "occured issue in processing script(docker_img_ls.sh)",
		}
	}

	return res, nil
}

func (s *ApiService) ImageRm(reqs *request.PostImageRmRequest) (map[string][]string, error) {
	hostLocalSepSlash := filepath.FromSlash(reqs.Host)
	hostIpAndPort, _ := strings.CutPrefix(hostLocalSepSlash, filepath.FromSlash("tcp://"))
	dockerDeamonHost := filepath.FromSlash(fmt.Sprintf("tcp://%s", hostIpAndPort))

	//dockerRepoUrl := fmt.Sprintf("%s/%s", s.Cfg.RestIp, s.Cfg.RestPort)

	executor := shell.NewScriptExecutor(s.Cfg.ShellDir)
	script := fmt.Sprintf(`
		echo %s | sudo -S docker -H %s image rm %s
	`, s.Cfg.BuildSvrPasswd, dockerDeamonHost, reqs.ImageName)

	result, err := executor.Execute(context.Background(), s.Cfg.ShellName, "docker_img_rm.sh", script, true)
	if err != nil {
		return nil, httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(docker_img_rm.sh): %v", err), response.FAIL)
	}

	logger.Println("<=== Run docker_img_rm.sh ===>")
	logger.Printf("content: %s", script)

	res := make(map[string][]string)

	logger.Println("=== Execution Result ===")
	res["execut_result"] = append(res["execut_result"], "=== Execution Result ===")
	logger.Printf("Duration: %v\n", result.Duration)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Duration: %v", result.Duration))
	logger.Printf("Exit Code: %d\n", result.ExitCode)
	res["execute_result"] = append(res["execute_result"], fmt.Sprintf("Exit Code: %d", result.ExitCode))

	logger.Println("=== STDOUT ===")
	res["stdout"] = append(res["stdout"], "=== STDOUT ===")
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

	if result.ExitCode != 0 {
		return res, &service.ShellError{
			ExitCode: result.ExitCode,
			Msg:      "occured issue in processing script(docker_img_rm.sh)",
		}
	}

	return res, nil
}
