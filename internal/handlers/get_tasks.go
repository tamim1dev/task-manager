package handlers

import (
	"net/http"
	"strconv"

	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("X-User-Id")
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")
	// defaults for pagination
	page := 1
	limit := 10

	if pageStr != "" {
		if parsedPage, pageErr := strconv.Atoi(pageStr); pageErr != nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr != "" {
		if parsedLimit, limitErr := strconv.Atoi(limitStr); limitErr != nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				limit = 100
			} else {
				limit = parsedLimit
			}
		}
	}

	offset := (page - 1) * limit

	tasks, dbErr := services.GetTasksByUserID(userId, limit, offset, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error on get tasks")
		return
	}
	utils.ReturnJson(w, http.StatusOK, tasks)
}
