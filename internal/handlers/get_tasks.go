package handlers

import (
	"net/http"

	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("X-User-Id")
	tasks, dbErr := services.GetTasksByUserID(userId, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error on get tasks")
		return
	}
	utils.ReturnJson(w, http.StatusOK, tasks)
}
