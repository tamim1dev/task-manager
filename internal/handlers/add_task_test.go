package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddTask(t *testing.T) {
	truncateTables(t)
	// add test user
	seedUserViaRegister(t, "abc", "abc@gmail.com", "abcpass")
	// login to get jwt
	token := getJwt(t, "abc@gmail.com", "abcpass")

	// main test
	payload := `{"title":"Buy groceries","description":"Milk, eggs, bread","days":7}`
	req := newAuthedRequest(http.MethodPost, "/tasks", payload, token)
	rec := httptest.NewRecorder()
	router := newTaskRouter()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		t.Errorf("expected %d, got %d, body: %s", http.StatusCreated, res.StatusCode, body)
	}
}
