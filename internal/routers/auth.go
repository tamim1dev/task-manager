package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/handlers"
	"github.com/tamim1dev/task-manager/internal/middleware"
)

func AuthRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	// register new user
	authRouter.Post("/register", handlers.RegisterUser)
	// login user
	authRouter.Post("/login", handlers.LoginUser)
	// jwt middleware check
	authRouter.Get("/me", middleware.AuthMiddleware(handlers.GetMe))

	return authRouter
}
