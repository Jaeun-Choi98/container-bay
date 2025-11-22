package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"net/http"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/Jaeun-Choi98/modules/shell"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v6"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
)

func (t *Controller) GetDockerPs(c *gin.Context) {
	ip := c.Param("ip")
	if ip == "" {
		c.Error(httperr.BADREQUEST.Add(nil, response.INVAILD_DATA))
		return
	}

	executor := shell.NewScriptExecutor("/tmp/scripts")

	script := fmt.Sprintf(`
		docker -H tcp://%s:2375 ps -a
	`, ip)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := executor.Execute(ctx, "bash", "docker_ps.sh", script, true)

	var res []string

	fmt.Println("=== Execution Result ===")
	res = append(res, "=== Execution Result ===")

	fmt.Printf("Duration: %v\n", result.Duration)
	res = append(res, fmt.Sprintf("Duration: %v", result.Duration))

	fmt.Printf("Exit Code: %d\n", result.ExitCode)
	res = append(res, fmt.Sprintf("Exit Code: %d", result.ExitCode))

	fmt.Println("=== STDOUT ===")
	res = append(res, "=== STDOUT ===")

	for _, line := range result.Stdout {
		fmt.Println(line)
		res = append(res, line)
	}
	fmt.Println("=== STDERR ===")
	res = append(res, "=== STDERR ===")
	for _, line := range result.Stderr {
		fmt.Println(line)
		res = append(res, line)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		res = append(res, fmt.Sprintf("Error: %v", err))
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}

func (t *Controller) PostBuildProject(c *gin.Context) {
	reqs := &request.PostBuildProjectRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	pjtPath := filepath.Join(t.Cfg.RepoDir, reqs.PjtName)
	log.Printf("[debug] pjt_path: %s", pjtPath)

	// clone하기 전에, 해당 path에 폴더가 이미 있는지 확인.
	os.RemoveAll(pjtPath)

	// repo는 필요하지 않음.
	_, err := git.PlainClone(pjtPath, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: t.Cfg.GitId,
			Password: t.Cfg.GitPwd,
		},
		URL:               reqs.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	buildPath := filepath.Join(pjtPath, reqs.ContextPath)
	log.Printf("[debug] build_path: %s", buildPath)

	c.JSON(http.StatusOK, "success")
}
