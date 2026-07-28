package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamim1dev/task-manager/internal/handlers"
)

func TestLoginUser(t *testing.T) {
	truncateTables(t)

	// first add user
	name := "abc"
	email := "abc@gmail.com"
	password := "abcpass"

	seedUserViaRegister(t, name, email, password)

	// now login with that user
	loginPayload := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
	loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginPayload))
	loginReq.Header.Set("Content-type", "application/json")
	loginRecorder := httptest.NewRecorder()
	handlers.LoginUser(loginRecorder, loginReq)

	loginResponse := loginRecorder.Result()
	defer func() {
		_ = loginResponse.Body.Close()
	}()

	if loginResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, loginResponse.StatusCode)
	}
}
