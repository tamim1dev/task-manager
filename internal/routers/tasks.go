package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/handlers"
	"github.com/tamim1dev/task-manager/internal/middleware"
)

func TasksRouter() *chi.Mux {
	tasksRouter := chi.NewRouter()

	tasksRouter.Post("/", middleware.AuthMiddleware(handlers.AddTask))
	tasksRouter.Get("/", middleware.AuthMiddleware(handlers.GetAllTasks))
	tasksRouter.Get("/{task_id}", middleware.AuthMiddleware(handlers.GetTaskById))
	tasksRouter.Patch("/{task_id}", middleware.AuthMiddleware(handlers.EditTask))
	tasksRouter.Delete("/{task_id}", middleware.AuthMiddleware(handlers.DeleteTaskById))

	return tasksRouter
}
