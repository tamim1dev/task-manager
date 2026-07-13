package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/handlers"
	"github.com/tamim1dev/task-manager/internal/middleware"
)

func UsersRouter() *chi.Mux {
	usersRouter := chi.NewRouter()

	usersRouter.Get("/me", middleware.AuthMiddleware(handlers.GetMe))

	return usersRouter
}
