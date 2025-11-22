package controller

import (
	"net/http"

	"github.com/Jaeun-Choi98/container-bay/internal/config"
	"github.com/Jaeun-Choi98/container-bay/internal/service"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/middleware"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	ApiSvc service.ApiServiceInterface
	Router *gin.Engine
	Cfg    *config.Config
}

func NewController(router *gin.Engine, apiSvc service.ApiServiceInterface, cfg *config.Config) *Controller {
	controller := &Controller{
		Router: router,
		ApiSvc: apiSvc,
		Cfg:    cfg,
	}
	controller.RoutingPath()
	return controller
}

func (t *Controller) RoutingPath() {
	t.Router.Use(middleware.LogMiddleware())
	t.Router.Use(middleware.ErrorMiddleware())
	t.Router.Use(middleware.NewCORSMiddleware([]string{"http://localhost:3000",
		"http://localhost:3001", "http://localhost:3002", "*"}))
	t.Router.Use(middleware.StoreMemberIdToContext())

	t.Router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	t.Router.GET("/docker/ps/:ip", t.GetDockerPs)

	t.Router.POST("/build", t.PostBuildProject)

	// SPA 미들웨어 동작 방식에 의해 제일 아래에 배치
	//t.Router.Use(middleware.SpaHandlerOther("/vworld", "vworld"))
	//t.Router.Use(middleware.SpaHandlerRoot("build", "index.html"))
}
