package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/utils"
)

type ReturnNewUser struct {
	Name string `json:"name"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.CreateUserRequest
	var returnUser ReturnNewUser

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		utils.ReturnError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	password_hash, err := utils.HashPassword(newUser.Password)
	if err != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING name`
	user_err := database.DB.Pool.QueryRow(r.Context(), query, newUser.Name, newUser.Email, password_hash).Scan(&returnUser.Name)
	if user_err != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	utils.ReturnJson(w, http.StatusOK, returnUser)
}
