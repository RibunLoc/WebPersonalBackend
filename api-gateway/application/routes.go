package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// healthcheck
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// auth proxy
	r.Route("/auth", func(rt chi.Router) {
		rt.Post("/login", a.AuthHandler.Login)       // gRPC -> auth
		rt.Post("/register", a.AuthHandler.Register) // HTTP -> auth
	})

	a.router = r
}
