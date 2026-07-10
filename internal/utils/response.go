package utils

import (
	"encoding/json"
	"net/http"
)

func ReturnJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func ReturnError(w http.ResponseWriter, status int, message string) {
	ReturnJson(w, status, map[string]string{"error": message})
}
