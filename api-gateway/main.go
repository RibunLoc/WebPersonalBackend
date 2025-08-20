package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/RibunLoc/WebPersonalBackend/api-gateway/application"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := application.LoadConfig()

	app, err := application.New(ctx, cfg)
	if err != nil {
		fmt.Println("failed to init app:", err)
		return
	}

	if err := app.Start(ctx); err != nil {
		fmt.Println("server error:", err)
	}
}
