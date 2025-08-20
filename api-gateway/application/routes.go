package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (a *App) loadRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173", "https://holoc.id.vn"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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

	r.Route("/contact", func(rt chi.Router) {
		rt.Post("/", a.ContactHandler.Submit)
	})

	a.router = r
}
