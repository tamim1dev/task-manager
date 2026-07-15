package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

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
