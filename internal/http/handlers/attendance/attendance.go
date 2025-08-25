package attendance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/priyanka-16/attendance-app-backEnd/internal/models"
	"github.com/priyanka-16/attendance-app-backEnd/internal/storage"
	"github.com/priyanka-16/attendance-app-backEnd/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var attendance models.Attendance

		err := json.NewDecoder(r.Body).Decode(&attendance)
		if errors.Is(err, io.EOF) { // error when no input is available
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err = validator.New().Struct(attendance); err != nil {
			validateErrs := err.(validator.ValidationErrors) //type cast required as simple errors cant be sent as argument to ValidationError
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateAttendance(
			&attendance,
		)
		slog.Info("attendance recorded successfully", slog.String("userId", fmt.Sprint(attendance.StudentID)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]uint{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a attendance", slog.String("id", id))

		intId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		attendance, err := storage.GetAttendanceById(uint(intId))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, attendance)
	}
}

func GetAttendancesList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		attendances, err := storage.GetAttendancesList()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, attendances)
	}
}
