package controller

import (
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
	t.Router.Use(middleware.NewCORSMiddleware([]string{"http://localhost:8090", "http://localhost:3000",
		"http://localhost:3001", "http://localhost:3002", "*"}))
	t.Router.Use(middleware.StoreMemberIdToContext())

	// == Docker ==
	t.Router.POST("/ps", t.PostDockerPs)
	t.Router.POST("/build", t.PostBuildProject)
	t.Router.POST("/run", t.PostRunProject)
	t.Router.POST("/stop", t.PostDockerStop)
	t.Router.POST("/restart", nil)
	t.Router.POST("/rm", t.PostDockerRemove)

	// == File Upload ==
	t.Router.POST("/upload-file", t.PostUploadFile)
	// path 아래에 폴더가 복사됨.
	t.Router.POST("/upload-targz", t.PostUploadTarGz)

	// 개발해야할 api

	// SPA 미들웨어 동작 방식에 의해 제일 아래에 배치
	//t.Router.Use(middleware.SpaHandlerOther("/vworld", "vworld"))
	t.Router.Use(middleware.SpaHandlerRoot("build", "index.html"))
}
