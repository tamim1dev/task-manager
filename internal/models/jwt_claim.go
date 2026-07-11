package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Token string
}
