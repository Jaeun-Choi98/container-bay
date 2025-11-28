package controller

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/Jaeun-Choi98/modules/shell"
	"github.com/gin-gonic/gin"
)

func (t *Controller) PostUploadFile(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("path is empty"), response.INVAILD_DATA))
		return
	}

	reqFile, err := c.FormFile("file")
	if err != nil {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("failed to upload file: %v", err), response.FAIL))
		return
	}

	source := filepath.Join(t.Cfg.VolumeDir, path)
	if err := c.SaveUploadedFile(reqFile, source); err != nil {
		c.Error(httperr.INNER_ERROR.Add(fmt.Errorf("failed to save file: %v", err), response.FAIL))
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add("success"))
}

func (t *Controller) PostUploadTarGz(c *gin.Context) {
	path := c.PostForm("path")
	if path == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("path is empty"), response.INVAILD_DATA))
		return
	}

	reqFile, err := c.FormFile("file")
	if err != nil {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("failed to upload file: %v", err), response.FAIL))
		return
	}

	dir := filepath.Join(t.Cfg.VolumeDir, path)
	source := filepath.Join(dir, reqFile.Filename)
	if err := c.SaveUploadedFile(reqFile, source); err != nil {
		c.Error(httperr.INNER_ERROR.Add(fmt.Errorf("failed to save file: %v", err), response.FAIL))
		return
	}

	execurator := shell.NewScriptExecutor(t.Cfg.ShellDir)
	script := fmt.Sprintf(`
		tar -xzvf %s -C %s
	`, source, dir)

	result, err := execurator.Execute(context.Background(), t.Cfg.ShellName, "targz_extract.sh", script, true)
	if err != nil {
		c.Error(httperr.INNER_ERROR.Add(fmt.Errorf("[API Service] failed to execute script(targz_extract.sh): %v", err), response.FAIL))
		return
	}

	logger.Println("<=== Run targz_extract.sh ===>")
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

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}
