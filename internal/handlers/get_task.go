package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

// GetTask godoc
// @Summary Gets a task by id
// @Description Gets a task by id requires valid jwt token
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /tasks/{task_id} [get]
func GetTaskById(w http.ResponseWriter, r *http.Request) {
	task_id := chi.URLParam(r, "task_id")
	user_id := r.Header.Get("X-User-Id")
	task, dbErr := services.GetTaskById(task_id, user_id, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, dbErr.Error())
		return
	}

	utils.ReturnJson(w, http.StatusOK, task)
}
