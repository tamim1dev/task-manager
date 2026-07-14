package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	task_id := chi.URLParam(r, "task_id")
	user_id := r.Header.Get("X-User-Id")
	task, dbErr := services.GetTaskById(task_id, user_id, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error on get task")
		return
	}

	utils.ReturnJson(w, http.StatusOK, task)
}
