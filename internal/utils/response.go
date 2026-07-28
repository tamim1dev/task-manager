package utils

import (
	"encoding/json"
	"net/http"

	"github.com/tamim1dev/task-manager/internal/models"
)

func ReturnJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func ReturnError(w http.ResponseWriter, status int, message string) {
	errorType := &models.ErrorResponse{
		Error: message,
	}
	ReturnJson(w, status, errorType)
}
