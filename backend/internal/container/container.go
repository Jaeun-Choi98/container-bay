package container

import (
	"github.com/Jaeun-Choi98/container-bay/internal/config"
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	apiservice "github.com/Jaeun-Choi98/container-bay/internal/service/api-service"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/controller"
	"github.com/gin-gonic/gin"
)

var container *Container

type Container struct {
	Config *config.Config
	Rest   *rest.RESTServer
}

func NewContainer() (*Container, error) {
	if container != nil {
		return container, nil
	}

	// 환경 변수
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	// 커스텀 로거
	customLogger, err := logger.NewCustomLogger("", config.LogFileMaxAge)
	if err != nil {
		return nil, err
	}
	logger.SetLogger(customLogger)

	// Rest
	apiSvc := apiservice.NewApiService()
	controller := controller.NewController(gin.New(), apiSvc)
	rest := rest.NewRestServer(config, controller)

	container = &Container{
		Config: config,
		Rest:   rest,
	}
	return container, nil
}
