package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyanka-16/attendance-app-backEnd/internal/models"
	"github.com/priyanka-16/attendance-app-backEnd/internal/storage"
	"github.com/priyanka-16/attendance-app-backEnd/internal/utils/response"
)

var jwtSecret = []byte("super-secret-key")

type LoginOTPRequest struct {
	Mobile string `json:"mobile"`
	OTP    string `json:"otp"`
}

func LoginWithOTP(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginOTPRequest

		if err := json.NewDecoder(r.Body).Decode(&req); errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		} else if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		fmt.Print(req.OTP)
		// Fetch latest OTP
		otp, err := store.GetUserOTPByMobile(req.Mobile)
		fmt.Print(otp.Code)
		if err != nil || otp.Code != req.OTP || time.Now().After(time.Unix(otp.ExpiresAt, 0)) {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid or expired OTP")))
			return
		}

		// Check if user exists or create
		user, err := store.GetUserByMobile(req.Mobile)
		if err != nil {
			newUser := models.User{Mobile: req.Mobile}
			_, err := store.CreateUser(&newUser)
			if err != nil {
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
				return
			}
			user = &newUser
		}

		// Generate JWT
		claims := jwt.MapClaims{
			"userId": user.ID,
			"mobile": user.Mobile,
			"exp":    time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{
			"token": tokenString,
		})
	}
}
