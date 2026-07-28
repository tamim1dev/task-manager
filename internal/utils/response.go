package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tamim1dev/task-manager/internal/models"
)

func ReturnJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func ReturnError(w http.ResponseWriter, status int, message string) {
	errorInstance := &models.ErrorResponse{
		Error: message,
	}
	ReturnJson(w, status, errorInstance)
}
