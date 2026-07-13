package routers

import (
	"github.com/go-chi/chi/v5"
)

func TasksRouter() *chi.Mux {
	tasksRouter := chi.NewRouter()

	return tasksRouter
}
