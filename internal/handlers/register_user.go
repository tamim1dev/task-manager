package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.CreateUserRequest
	var returnUser models.ReturnNewUser

	jsonerr := json.NewDecoder(r.Body).Decode(&newUser)
	if jsonerr != nil {
		utils.ReturnError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// input validation
	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		utils.ReturnError(w, http.StatusBadRequest, "Name, email, and password are required")
		return
	}

	password_hash, hasherr := utils.HashPassword(newUser.Password)
	if hasherr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING name`
	user_err := database.DB.Pool.QueryRow(r.Context(), query, newUser.Name, newUser.Email, password_hash).Scan(&returnUser.Name)
	if user_err != nil {
		var pgErr *pgconn.PgError
		if errors.As(user_err, &pgErr) && pgErr.Code == "23505" {
			utils.ReturnError(w, http.StatusConflict, "Email already exists")
			return
		}
		utils.ReturnError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	utils.ReturnJson(w, http.StatusOK, returnUser)
}
