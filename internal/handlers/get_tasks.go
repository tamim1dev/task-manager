package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("X-User-Id")
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")
	// defaults for pagination
	var offset int
	limit := 10

	if parsedLimit, limitErr := strconv.Atoi(limitStr); limitErr == nil && parsedLimit > 0 {
		limit = min(parsedLimit, 10)
	}
	if parsedPage, pageErr := strconv.Atoi(pageStr); pageErr == nil && parsedPage > 0 {
		offset = (parsedPage - 1) * limit
	}

	// completed filter
	var completed *bool
	completedStr := query.Get("completed")
	if completedStr != "" {
		if parsedCompleted, completedErr := strconv.ParseBool(completedStr); completedErr == nil {
			completed = &parsedCompleted
		}
	}

	// sort_by and sort_order
	sortBy := "created_at"
	if query.Get("sort_by") == "due_date" {
		sortBy = "due_date"
	}
	sortOrder := "DESC"
	if strings.ToLower(query.Get("sort_order")) == "asc" {
		sortOrder = "ASC"
	}

	tasks, dbErr := services.GetTasksByUserID(userId, limit, offset, completed, sortBy, sortOrder, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error on get tasks")
		return
	}
	utils.ReturnJson(w, http.StatusOK, tasks)
}
