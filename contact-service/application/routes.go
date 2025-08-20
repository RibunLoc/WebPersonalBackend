package application

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (a *App) loadRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	ch := a.buildHandlers()
	r.Post("/contact", ch.Submit)
	r.Post("/contact/dev/sendmail", ch.DevSendMail)

	a.router = r
}
