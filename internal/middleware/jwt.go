package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tamim1dev/task-manager/internal/models"
	"github.com/tamim1dev/task-manager/internal/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		autheHeader := r.Header.Get("Authorization")
		if autheHeader == "" {
			utils.ReturnError(w, http.StatusUnauthorized, "Token required")
			return
		}

		bearerSplit := strings.Split(autheHeader, " ")
		if len(bearerSplit) != 2 || bearerSplit[0] != "Bearer" {
			utils.ReturnError(w, http.StatusUnauthorized, "Invalid bearer format")
			return
		}

		tokenString := bearerSplit[1]
		claims := &models.JwtClaims{}

		token, tokenErr := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if tokenErr != nil || !token.Valid {
			utils.ReturnError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		r.Header.Set("X-User-Email", claims.Email)
		next.ServeHTTP(w, r)
	}
}
