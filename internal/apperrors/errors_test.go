package apperrors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	type fields struct {
		Facade string
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "All arguments",
			fields: fields{
				Facade: "Facade",
				Err:    errors.New("error"),
			},
			want: "[Facade] error",
		},
		{
			name:   "Facade only",
			fields: fields{Facade: "Facade"},
			want:   "Facade",
		},
		{
			name:   "Error only",
			fields: fields{Err: errors.New("error")},
			want:   "error",
		},
		{
			name: "No arguments",
			want: "[Empty]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := AppError{
				Facade: tt.fields.Facade,
				Err:    tt.fields.Err,
			}

			assert.Equal(t, tt.want, e.Error())
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	err := errors.New("error")

	type fields struct {
		Facade string
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "All arguments",
			fields: fields{
				Facade: "Facade",
				Err:    err,
			},
			want: err,
		},
		{
			name:   "Facade only",
			fields: fields{Facade: "Facade"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := AppError{
				Facade: tt.fields.Facade,
				Err:    tt.fields.Err,
			}

			assert.Equal(t, tt.want, e.Unwrap())
		})
	}
}

func TestEmptyError(t *testing.T) {
	tests := []struct {
		name string
		want *AppError
	}{
		{
			name: "Empty error",
			want: &AppError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, EmptyError())
		})
	}
}

func TestNewError(t *testing.T) {
	err := errors.New("error")

	type args struct {
		text string
		err  error
	}
	tests := []struct {
		name string
		args args
		want *AppError
	}{
		{
			name: "All arguments",
			args: args{
				text: "Facade",
				err:  err,
			},
			want: &AppError{
				Facade: "Facade",
				Err:    err,
			},
		},
		{
			name: "Facade only",
			args: args{text: "Facade"},
			want: &AppError{Facade: "Facade"},
		},
		{
			name: "Error only",
			args: args{err: err},
			want: &AppError{Err: err},
		},
		{
			name: "No arguments",
			want: &AppError{Facade: "", Err: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewError(tt.args.text, tt.args.err))
		})
	}
}

func TestHandleHTTPError(t *testing.T) {
	type args struct {
		err  *AppError
		code int
	}
	tests := []struct {
		name string
		args args
		want *AppError
	}{
		{
			name: "Existing error",
			args: args{
				err:  NewError("err", errors.New("err")),
				code: 500,
			},
			want: NewError("", errors.New("err")),
		},
		{
			name: "Missing facade",
			args: args{
				err:  nil,
				code: 500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleHTTPError(w, tt.args.err, tt.args.code)
			res := w.Result()
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			assert.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
			assert.Equal(t, tt.args.code, w.Result().StatusCode)
		})
	}
}

func TestHandleInternalError(t *testing.T) {
	tests := []struct{ name string }{{name: "Testing function"}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleInternalError(w)
			res := w.Result()
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			assert.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		})
	}
}

func TestHandleURLError(t *testing.T) {
	tests := []struct{ name string }{{name: "Testing function"}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleURLError(w)
			res := w.Result()
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			assert.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		})
	}
}

func TestHandleUserError(t *testing.T) {
	tests := []struct{ name string }{{name: "Testing function"}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleUserError(w)
			res := w.Result()
			assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			assert.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		})
	}
}
