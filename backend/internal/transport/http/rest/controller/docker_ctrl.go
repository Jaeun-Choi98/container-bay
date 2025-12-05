package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Jaeun-Choi98/container-bay/internal/service"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/gin-gonic/gin"
)

func (t *Controller) PostDockerPs(c *gin.Context) {
	reqs := &request.PostDockerPsRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	res, err := t.ApiSvc.DockerPs(reqs.Host)
	if err != nil {
		var shellError *service.ShellError
		if errors.As(err, &shellError) {
			c.JSON(http.StatusOK, response.FAIL.Add(res))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}

func (t *Controller) PostBuildProject(c *gin.Context) {
	reqs := &request.PostBuildProjectRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.ContextPath == "" || reqs.PjtName == "" || reqs.URL == "" {
		c.Error(httperr.BADREQUEST.Add(nil, response.INVAILD_DATA))
		return
	}

	res, err := t.ApiSvc.CloneAndBuild(reqs.URL, reqs.PjtName, reqs.ContextPath)
	if err != nil {
		var shellError *service.ShellError
		if errors.As(err, &shellError) {
			c.JSON(http.StatusOK, response.FAIL.Add(res))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}

func (t *Controller) PostRunProject(c *gin.Context) {
	reqs := &request.PostRunProjectRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.Host == "" || reqs.Name == "" || reqs.Image == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("need info(host, name, image)"), response.INVAILD_DATA))
		return
	}

	res, err := t.ApiSvc.RunContainer(reqs)
	if err != nil {
		var shellError *service.ShellError
		if errors.As(err, &shellError) {
			c.JSON(http.StatusOK, response.FAIL.Add(res))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}

func (t *Controller) PostDockerStop(c *gin.Context) {
	reqs := &request.PostDockerStopRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.ContainerName == "" || reqs.Host == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("need info(container name, docker daemon host)"), response.INVAILD_DATA))
		return
	}

	res, err := t.ApiSvc.StopContainer(reqs)
	if err != nil {
		var shellError *service.ShellError
		if errors.As(err, &shellError) {
			c.JSON(http.StatusOK, response.FAIL.Add(res))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}

func (t *Controller) PostDockerRemove(c *gin.Context) {
	reqs := &request.PostDockerRemoveRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.ContainerName == "" || reqs.Host == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("need info(container name, docker daemon host)"), response.INVAILD_DATA))
		return
	}

	res, err := t.ApiSvc.RemoveContainer(reqs)
	if err != nil {
		var shellError *service.ShellError
		if errors.As(err, &shellError) {
			c.JSON(http.StatusOK, response.FAIL.Add(res))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}

func (t *Controller) PostDockerRestart(c *gin.Context) {
	reqs := &request.PostDockerRestartRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.ContainerName == "" || reqs.Host == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("need info(container name, docker daemon host)"), response.INVAILD_DATA))
		return
	}

	res, err := t.ApiSvc.RestartContainer(reqs)
	if err != nil {
		var shellError *service.ShellError
		if errors.As(err, &shellError) {
			c.JSON(http.StatusOK, response.FAIL.Add(res))
			return
		}
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}
