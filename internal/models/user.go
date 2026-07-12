package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password_Hash string    `json:"-"`
	Created_At    time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReturnUserInfo struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Created_At time.Time `json:"created_at"`
}
