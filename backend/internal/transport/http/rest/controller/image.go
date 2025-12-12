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

func (t *Controller) PostImageLs(c *gin.Context) {
	reqs := &request.PostImageLsRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.Host == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("need info(docker daemon host)"), response.FAIL))
		return
	}

	res, err := t.ApiSvc.ImageLs(reqs.Host)
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

func (t *Controller) PostImageRm(c *gin.Context) {
	reqs := &request.PostImageRmRequest{}
	if err := c.BindJSON(reqs); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}

	if reqs.Host == "" || reqs.ImageName == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("need info(image name, docker daemon host)"), response.FAIL))
		return
	}

	res, err := t.ApiSvc.ImageRm(reqs)
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
