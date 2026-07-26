package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginUser(t *testing.T) {
	truncateTables(t)

	// first add user
	name := "abc"
	email := "abc@gmail.com"
	password := "abcpass"

	registerPayload := fmt.Sprintf(`{"name":%q,"email":%q,"password":%q}`, name, email, password)
	registerReq := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(registerPayload))
	registerReq.Header.Set("Content-type", "application/json")
	registerRecorder := httptest.NewRecorder()
	RegisterUser(registerRecorder, registerReq)

	registerResponse := registerRecorder.Result()
	defer registerResponse.Body.Close()

	if registerResponse.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, registerResponse.StatusCode)
	}

	// now login with that user
	loginPayload := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
	loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginPayload))
	loginReq.Header.Set("Content-type", "application/json")
	loginRecorder := httptest.NewRecorder()
	LoginUser(loginRecorder, loginReq)

	loginResponse := loginRecorder.Result()
	defer loginResponse.Body.Close()

	if loginResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, loginResponse.StatusCode)
	}
}
