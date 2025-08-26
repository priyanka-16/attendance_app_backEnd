package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/priyanka-16/attendance-app-backEnd/internal/models"
	"github.com/priyanka-16/attendance-app-backEnd/internal/storage"
	"github.com/priyanka-16/attendance-app-backEnd/internal/utils/response"
)

type RequestOTPRequest struct {
	Mobile string `json:"mobile" validate:"required,len=10"`
}

func RequestOTP(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestOTPRequest

		if err := json.NewDecoder(r.Body).Decode(&req); errors.Is(err, io.EOF) { // error when no input is available
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		} else if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate input
		if err := validator.New().Struct(req); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Generate OTP (simple numeric for now)
		otpCode := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

		otp := models.UserOTP{
			Mobile:    req.Mobile,
			Code:      otpCode,
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		}

		_, err := store.CreateUserOTP(&otp)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// TODO: send OTP via SMS gateway
		response.WriteJson(w, http.StatusOK, map[string]string{
			"message": "OTP generated successfully",
			"otp":     otpCode, // ⚠️ In prod, don’t return OTP in response
		})
	}
}
