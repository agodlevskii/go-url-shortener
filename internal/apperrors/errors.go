package apperrors

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

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

type AppError struct {
	Facade string
	Err    error
}

func NewError(text string, err error) *AppError {
	return &AppError{
		Facade: text,
		Err:    err,
	}
}

func EmptyError() *AppError {
	return &AppError{}
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Facade, e.Err)
}

func (e AppError) Unwrap() error {
	return e.Err
}

func HandleHTTPError(w http.ResponseWriter, err *AppError, code int) {
	if err.Facade == "" {
		err.Facade = http.StatusText(code)
	}

	if err.Err != nil {
		log.Error(err.Error())
	}

	http.Error(w, err.Facade, code)
}

func HandleInternalError(w http.ResponseWriter) {
	HandleHTTPError(w, EmptyError(), http.StatusInternalServerError)
}

func HandleURLError(w http.ResponseWriter) {
	HandleHTTPError(w, NewError(URLFormat, nil), http.StatusBadRequest)
}

func HandleUserError(w http.ResponseWriter) {
	HandleHTTPError(w, NewError(UserID, nil), http.StatusBadRequest)
}
