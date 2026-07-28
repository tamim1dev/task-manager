package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTask(t *testing.T) {
	truncateTables(t)

	user := setupAuthedUser(t)
	taskID := setupTaskForUser(t, user.Token)

	router := newTaskRouter()

	req := newAuthedRequest(http.MethodGet, "/tasks/"+taskID, "", user.Token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("expected %d, got %d, body: %s", http.StatusOK, res.StatusCode, body)
	}
}
