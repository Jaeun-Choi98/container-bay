package main

import (
	"context"
	"log"

	"github.com/Jaeun-Choi98/container-bay/internal/app"
	"github.com/Jaeun-Choi98/container-bay/internal/container"
)

func main() {
	mycontainer, err := container.NewContainer()
	if err != nil {
		log.Printf("[Main] failed to create container: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	app := app.NewApplication(mycontainer, cancel)

	app.Start()
	<-ctx.Done()
}
