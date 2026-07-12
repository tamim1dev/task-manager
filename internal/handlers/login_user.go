package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginUserRequest

	jsonerr := json.NewDecoder(r.Body).Decode(&loginRequest)
	if jsonerr != nil {
		utils.ReturnError(w, http.StatusBadRequest, "Invalid json")
		return
	}
	defer r.Body.Close()

	if loginRequest.Email == "" || loginRequest.Password == "" {
		utils.ReturnError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	userFromDb, userGetErr := services.GetUserByEmail(loginRequest.Email, r)
	if userGetErr != nil {
		utils.ReturnError(w, http.StatusNotFound, "User with the email not found")
		return
	}

	passwordMatched := utils.CheckPasswordHash(loginRequest.Password, userFromDb.Password_Hash)
	if !passwordMatched {
		utils.ReturnError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// jwt
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.JwtClaims{
		Email: loginRequest.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, tokenErr := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if tokenErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	tokenResponse := &models.TokenResponse{
		Token: tokenStr,
	}
	utils.ReturnJson(w, http.StatusOK, tokenResponse)
}
