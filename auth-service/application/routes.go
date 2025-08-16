package application

import (
	"net/http"

	"github.com/RibunLoc/microservices-learn/auth-service/handler"
	"github.com/RibunLoc/microservices-learn/auth-service/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/auth", a.loadUserLogin)

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {

	userHandler := &handler.UserRegister{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
		},
	}

	router.Post("/", userHandler.CreateUserHandler)
}

func (a *App) loadUserLogin(router chi.Router) {
	userHandler := &handler.UserLogin{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
			JwtSecret:  a.config.JwtSecret,
		},
	}

	userRegiterHandler := &handler.UserRegister{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
		},
	}

	router.Post("/register", userRegiterHandler.CreateUserHandler)
	router.Post("/login", userHandler.LoginHandler)
	//router.Post("/logout", userHandler.LogoutHandler)
	//router.Post("/refresh-token", userHandler.RefreshTokenHandler)
	//router.Get("/me", userHandler.GetMeHandler)
}
