package handlers

import (
	"net/http"

	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/services"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func GetMe(w http.ResponseWriter, r *http.Request) {
	userMail := r.Header.Get("X-User-Email")
	userFromDb, dbErr := services.GetUserByEmail(userMail, r)
	if dbErr != nil {
		utils.ReturnError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	returnUser := &models.ReturnUserInfo{
		Id:         userFromDb.Id,
		Name:       userFromDb.Name,
		Email:      userFromDb.Email,
		Created_At: userFromDb.Created_At,
	}
	utils.ReturnJson(w, http.StatusOK, returnUser)
}
