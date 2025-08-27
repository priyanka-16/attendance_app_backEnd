package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/priyanka-16/attendance-app-backEnd/internal/models"
	"github.com/priyanka-16/attendance-app-backEnd/internal/storage"
	"github.com/priyanka-16/attendance-app-backEnd/internal/utils/response"
)

// JWT secret (move this to config later)
var jwtSecret = []byte("super-secret-key")

// Claims defines JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GetProfile handler
func GetProfile(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT
		userID, err := extractUserIDFromToken(r)
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
			return
		}

		// Fetch user
		// user, err := store.GetUserById(userID)
		// if err != nil {
		// 	response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		// 	return
		// }

		// Fetch student profile (example)
		student, err := store.GetUserStudentById(userID)
		if err != nil {
			// fallback: maybe user is a teacher, or return basic user info
			slog.Info("user profile not found in student table", slog.Int64("id", int64(userID)))
		}

		profile := map[string]interface{}{
			"student": student,
		}

		response.WriteJson(w, http.StatusOK, profile)
	}
}

// UpdateProfile handler
func UpdateProfile(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID
		userID, err := extractUserIDFromToken(r)
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
			return
		}

		// Decode request body
		var payload struct {
			Name     string `json:"name" validate:"required"`
			Photo    string `json:"photo"`
			State    string `json:"state"`
			District string `json:"district"`
		}

		err = json.NewDecoder(r.Body).Decode(&payload)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate
		if err = validator.New().Struct(payload); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		student := models.UserStudent{
			UserID:   uint(userID),
			Name:     payload.Name,
			Photo:    payload.Photo,
			State:    payload.State,
			District: payload.District,
			// You’d add Photo, State, District fields to UserStudent if needed
		}

		err = store.UpdateUserStudent(student)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "profile updated"})
	}
}

// Helper: extract userID from JWT
func extractUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing Authorization header")
	}

	tokenString := authHeader[len("Bearer "):]
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
