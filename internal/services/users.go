package services

import (
	"net/http"

	"github.com/tamim1dev/task-manager/internal/database"
	"github.com/tamim1dev/task-manager/internal/models"
)

func GetUserByEmail(email string, r *http.Request) (models.User, error) {
	var targetUser models.User
	query := `SELECT * FROM users WHERE email = ($1)`
	dbErr := database.DB.Pool.QueryRow(r.Context(), query, email).Scan(&targetUser)
	if dbErr != nil {
		return models.User{}, dbErr
	}
	return targetUser, nil
}
