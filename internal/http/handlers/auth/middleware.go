package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyanka-16/attendance-app-backEnd/internal/utils/response"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}

		next(w, r)
	}
}
