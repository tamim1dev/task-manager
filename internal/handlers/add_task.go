package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("X-User-Id")
	var taskRequest models.AddTaskRequest

	jsonErr := json.NewDecoder(r.Body).Decode(&taskRequest)
	if jsonErr != nil {
		utils.ReturnError(w, http.StatusBadRequest, "title, description, days required")
		return
	}

	newTask := &models.AddTask{
		Title:       taskRequest.Title,
		Description: taskRequest.Description,
		Due_Date:    time.Now().Add(time.Duration(taskRequest.Days) * time.Hour * 24),
		User_Id:     uuid.MustParse(userId),
	}

	addedTask, dbErr := services.AddTask(newTask, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Can't save task")
		return
	}

	utils.ReturnJson(w, http.StatusOK, addedTask)
}
