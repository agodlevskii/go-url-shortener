// Package middlewares provides custom middleware for the chi router, and all related functionality.
package middlewares

import (
	"crypto/aes"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/encryptors"
)

type AuthConfig interface {
	GetUserCookieName() string
}

// Authorize provides a cookie-based user authorization.
// If the cookie is present and valid, Authorize passes the execution to the next handler.
// If the cookie is missing or invalid, Authorize creates a new cookie and adds it to the request.
func Authorize(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cfg.GetUserCookieName())

			if err == nil {
				valid, vErr := validateID(cookie.Value)
				if vErr != nil {
					apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
					return
				}

				if valid {
					next.ServeHTTP(w, r)
					return
				}
			}

			newID, err := generateID()
			if err != nil {
				apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
				return
			}

			cookie = &http.Cookie{Name: cfg.GetUserCookieName(), Value: newID, Path: "/"}
			http.SetCookie(w, cookie)
			r.AddCookie(cookie)
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID parses the request's user-related cookie and decrypts its value via encryptors.AESDecrypt functionality.
// If the cookie is missing, or its value fails to be decrypted, the error will be returned.
func GetUserID(cfg AuthConfig, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cfg.GetUserCookieName())
	if err != nil {
		return "", err
	}

	id, err := encryptors.AESDecrypt(cookie.Value)
	if err != nil {
		return "", err
	}
	return string(id), err
}

// generateID generates new encrypted ID in the UUID format.
// The encryption is performed by the encryptors.AESEncrypt functionality.
func generateID() (string, error) {
	id := uuid.New()
	return encryptors.AESEncrypt(id.String()[:aes.BlockSize])
}

// validateID decrypts the ID value and checks if it matches the UUID format.
// If the decryption fails, the false value and the error will be returned.
func validateID(id string) (bool, error) {
	data, err := encryptors.AESDecrypt(id)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return regexp.Match(`\w{8}-\w{4}-\w{2}`, data)
}
