package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status	string		`json:"status"`
	Error	string		`json:"error,omitempty"`
}

const(
	StatusOK = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error: err,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs, ", "),
	}
}