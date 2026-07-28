package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamim1dev/task-manager/internal/handlers"
)

func TestRegisterUser(t *testing.T) {
	truncateTables(t)
	payload := `{"name":"abc","email":"abc@gmail.com","password":"abcpass"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-type", "application/json")
	recorder := httptest.NewRecorder()

	handlers.RegisterUser(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, resp.StatusCode)
	}
}
