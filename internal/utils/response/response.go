package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response { //All errors are received as argument
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field())) //Sprintf is used to concatenate strings
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field())) //Sprintf is used to concatenate strings

		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, " ,"), //strings.Joins helps to concatenate slised into strings. 1st arg- slice 2nd arg - separator
	}
}
