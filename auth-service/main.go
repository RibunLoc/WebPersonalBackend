package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/RibunLoc/microservices-learn/auth-service/application"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	app, err := application.New(ctx, application.LoadConfig())

	if err != nil {
		fmt.Println("failed to load config server: ", err)
		return
	}

	defer cancel()

	errStart := app.Start(ctx)
	if errStart != nil {
		fmt.Println("failed to start app: ", err)
	}
	cancel()
}
