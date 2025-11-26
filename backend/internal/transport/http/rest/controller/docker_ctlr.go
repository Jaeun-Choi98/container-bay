package controller

import (
	"fmt"
	"net/http"

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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS.Add(res))
}
