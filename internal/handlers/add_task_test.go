package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamim1dev/task-manager/internal/middleware"
)

func TestAddTask(t *testing.T) {
	truncateTables(t)
	// add test user
	seedUserViaRegister(t, "abc", "abc@gmail.com", "abcpass")
	// login to get jwt
	token := getJwt(t, "abc@gmail.com", "abcpass")

	// main test
	payload := `{"title":"Buy groceries","description":"Milk, eggs, bread","days":7}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler := middleware.AuthMiddleware(AddTask)
	handler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Errorf("expected %d, got %d, body: %s", http.StatusCreated, res.StatusCode, body)
	}
}
