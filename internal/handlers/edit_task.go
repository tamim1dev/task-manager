package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

// UpdateTask godoc
// @Summary Updates a task by id
// @Description Updates a task by id requires valid jwt token
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /tasks/{task_id} [patch]
func EditTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "task_id")
	userId := r.Header.Get("X-User-Id")
	var changes models.UpdateTask

	jsonErr := json.NewDecoder(r.Body).Decode(&changes)
	if jsonErr != nil {
		utils.ReturnError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	editedTask, editErr := services.EditTask(changes, taskId, userId, r)
	if editErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, editErr.Error())
		return
	}

	utils.ReturnJson(w, http.StatusOK, editedTask)
}
