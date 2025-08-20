package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RibunLoc/WebPersonalBackend/api-gateway/handler"
	"github.com/RibunLoc/WebPersonalBackend/api-gateway/internal/client"
)

type App struct {
	cfg            Config
	router         http.Handler
	AuthHandler    *handler.AuthProxy
	ContactHandler *handler.ContactProxy

	// closers
	closeAuthGRPC    func() error
	closeContactGRPC func() error
}

func New(ctx context.Context, cfg Config) (*App, error) {
	authGRPC, closeFn, err := client.NewAuthGRPC(cfg.AuthGRPCAddr)
	if err != nil {
		return nil, fmt.Errorf("connect auth gRPC: %w", err)
	}

	contactGRPC, closeContact, err := client.NewContactGRPC(cfg.ContactGRPCAddr)
	if err != nil {
		return nil, fmt.Errorf("connect contact gRPC: %w", err)
	}

	app := &App{
		cfg:              cfg,
		AuthHandler:      handler.NewAuthProxy(authGRPC, cfg.AuthHTTPBase),
		ContactHandler:   handler.NewContactProxy(contactGRPC),
		closeAuthGRPC:    closeFn,
		closeContactGRPC: closeContact,
	}
	app.loadRoutes()
	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.ServerPort),
		Handler: a.router,
	}
	errCh := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		_ = server.Shutdown(ctx)
		_ = a.cleanup()
		return nil
	}
}

func (a *App) cleanup() error {
	if a.closeAuthGRPC != nil {
		_ = a.closeAuthGRPC()
	}
	if a.closeContactGRPC != nil {
		_ = a.closeContactGRPC()
	}
	return nil
}
