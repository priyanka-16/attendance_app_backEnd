package school_grade_section

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
		var section models.SchoolGradeSection

		err := json.NewDecoder(r.Body).Decode(&section)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err = validator.New().Struct(section); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateSchoolGradeSection(&section)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("school grade section created successfully", slog.String("schoolGradeSectionId", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]uint{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting school grade section", slog.String("id", id))

		intId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		section, err := storage.GetSchoolGradeSectionById(uint(intId))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, section)
	}
}

func GetSchoolGradeSectionsList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sections, err := storage.GetSchoolGradeSectionsList()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, sections)
	}
}
