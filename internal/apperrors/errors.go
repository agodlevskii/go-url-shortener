// Package apperrors includes the custom application errors.
package apperrors

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// The constants list all possible custom error messages.
const (
	URLGone          = "the requested URL is no longer available"
	URLFormat        = "you provided an incorrect URL format"
	URLNotFound      = "the requested URL not found"
	UserID           = "cannot identify the user"
	BatchFormat      = "you provided an incorrect batch format"
	IDsListFormat    = "you provided an incorrect IDs list format"
	IDSize           = "the ID size is missing"
	IDGeneration     = "cannot generate the ID"
	RandomStrLen     = "random string length is missing"
	FilenameMissing  = "the filename is missing"
	FileMalformed    = "the file is malformed"
	RepoEntryInvalid = "the stored entry is invalid"
)

// AppError describes a custom error.
// Facade includes a custom message for the existing error.
// Err represents the original error wrapped in the custom one.
type AppError struct {
	Facade string
	Err    error
}

// NewError returns a new instance of AppError.
func NewError(text string, err error) *AppError {
	return &AppError{
		Facade: text,
		Err:    err,
	}
}

// EmptyError returns a new instance of AppErrors with no detailed information about it.
// It's used in the function when the application returns a user error that's not based on the real application error.
func EmptyError() *AppError {
	return &AppError{}
}

// Error implements the error interface.
// If the error is empty, an empty string will be returned.
// In case if either facade or error is missing, the existing part will be returned.
func (e AppError) Error() string {
	if e.Facade == "" && e.Err == nil {
		return "[Empty]"
	}

	if e.Facade == "" {
		return e.Err.Error()
	}

	if e.Err == nil {
		return e.Facade
	}

	return fmt.Sprintf("[%s] %v", e.Facade, e.Err)
}

// Unwrap returns the original error wrapped into a custom one.
func (e AppError) Unwrap() error {
	return e.Err
}

// HandleHTTPError creates http.Error based on the custom AppError.
func HandleHTTPError(w http.ResponseWriter, err *AppError, code int) {
	if err.Facade == "" {
		err.Facade = http.StatusText(code)
	}

	if err.Err != nil {
		log.Error(err.Error())
	}

	http.Error(w, err.Facade, code)
}

// HandleInternalError fires the 500 HTTP status.
func HandleInternalError(w http.ResponseWriter) {
	HandleHTTPError(w, EmptyError(), http.StatusInternalServerError)
}

// HandleURLError handles the incorrect URLs provided by the user.
func HandleURLError(w http.ResponseWriter) {
	HandleHTTPError(w, NewError(URLFormat, nil), http.StatusBadRequest)
}

// HandleUserError handles the user-related errors, e.g. missing cookie value.
func HandleUserError(w http.ResponseWriter) {
	HandleHTTPError(w, NewError(UserID, nil), http.StatusBadRequest)
}
