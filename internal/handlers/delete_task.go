package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

// DeleteTask godoc
// @Summary Deletes a task by id
// @Description Deletes a task by id requires valid jwt token
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /tasks/{task_id} [delete]
func DeleteTaskById(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "task_id")
	userId := r.Header.Get("X-User-Id")
	deletedId, dbErr := services.DeleteTaskById(taskId, userId, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, dbErr.Error())
		return
	}

	utils.ReturnJson(w, http.StatusOK, deletedId)
}
