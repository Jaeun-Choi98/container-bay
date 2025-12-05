package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/container"
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/redis"
)

type Application struct {
	myContainer   *container.Container
	wg            sync.WaitGroup
	mainCtxCacnel context.CancelFunc
}

func NewApplication(c *container.Container, cancel context.CancelFunc) *Application {
	new := &Application{
		myContainer:   c,
		mainCtxCacnel: cancel,
	}
	return new
}

func (a *Application) Start() {
	// graceful shutdown
	detectSig := make(chan os.Signal, 1)
	signal.Notify(detectSig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-detectSig
		logger.Println("[App] Shutting down...")
		a.Shutdown()
		close(detectSig)
	}()

	// rest
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.myContainer.Rest.Start(); err != nil {
			logger.Printf("[App] Failed to start rest server: %v", err)
		}
	}()

	// log
	go logger.StartCleaning()
}

func (a *Application) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.myContainer.Rest.Shutdown(ctx); err != nil {
		logger.Printf("[App] issue in shutting down rest: %v", err)
	}
	a.wg.Wait()
	if err := redis.CloseRedisClient(); err != nil {
		logger.Printf("[App] issue in closing close redis client: %v", err)
	}
	logger.Println("[App] Application is terminated")
	logger.Shutdown()
	a.mainCtxCacnel()
}
