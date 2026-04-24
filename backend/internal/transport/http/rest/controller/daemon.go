package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/http-utils/httperr"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/response"
	"github.com/gin-gonic/gin"
)

func (t *Controller) GetDaemons(c *gin.Context) {
	daemons, err := t.ApiSvc.GetDaemons()
	if err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}
	c.JSON(http.StatusOK, response.SUCCESS.Add(daemons))
}

func (t *Controller) PostAddDaemon(c *gin.Context) {
	req := &request.PostAddDaemonRequest{}
	if err := c.BindJSON(req); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}
	if req.Host == "" {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("host is required"), response.INVAILD_DATA))
		return
	}
	if err := t.ApiSvc.AddDaemon(req); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}
	c.JSON(http.StatusOK, response.SUCCESS)
}

func (t *Controller) DeleteDaemon(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.Error(httperr.BADREQUEST.Add(fmt.Errorf("invalid daemon id"), response.INVAILD_DATA))
		return
	}
	if err := t.ApiSvc.RemoveDaemon(id); err != nil {
		c.Error(httperr.INNER_ERROR.Add(err, response.FAIL))
		return
	}
	c.JSON(http.StatusOK, response.SUCCESS)
}
