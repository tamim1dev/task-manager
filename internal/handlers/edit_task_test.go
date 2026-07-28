package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEditTask(t *testing.T) {
	truncateTables(t)

	user := setupAuthedUser(t)
	taskID := setupTaskForUser(t, user.Token)
	router := newTaskRouter()

	payload := `{"title":"Buy groceries and butter"}`
	req := newAuthedRequest(http.MethodPatch, "/tasks/"+taskID, payload, user.Token) // note: PATCH, matches your router
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("expected %d, got %d, body: %s", http.StatusOK, res.StatusCode, body)
	}
}
