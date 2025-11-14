package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/Jaeun-Choi98/modules/shell"
	"github.com/gin-gonic/gin"
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

func (t *Controller) PostDockerfileBuild(c *gin.Context) {

}
