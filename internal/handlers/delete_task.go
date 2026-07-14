package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

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
