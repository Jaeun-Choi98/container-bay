package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/config"
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/controller"
)

type RESTServer struct {
	server *http.Server
}

func NewRestServer(cfg *config.Config, controller *controller.Controller) *RESTServer {
	svr := &http.Server{
		Addr: cfg.RestIp + ":" + cfg.RestPort,
		// if using sse, comment below
		// WriteTimeout: time.Second * 5,
		ReadTimeout: time.Second * 5,
		Handler:     controller.Router,
	}
	return &RESTServer{
		server: svr,
	}
}

func (r *RESTServer) Start() error {
	return r.server.ListenAndServe()
}

func (r *RESTServer) Shutdown(ctx context.Context) (err error) {
	err = r.server.Shutdown(ctx)
	logger.Println("[REST] REST goroutine terminated")
	return err
}
