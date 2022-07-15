package middlewares

import (
	"crypto/aes"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/configs"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/encryptors"
	"net/http"
	"regexp"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(configs.UserCookieName)

		if err == nil {
			valid, err := validateID(cookie.Value)
			if err != nil {
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

		cookie = &http.Cookie{Name: configs.UserCookieName, Value: newID, Path: "/"}
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(configs.UserCookieName)
	if err != nil {
		return "", err
	}

	id, err := encryptors.AESDecrypt(cookie.Value)
	if err != nil {
		return "", err
	}
	return string(id), err
}

func generateID() (string, error) {
	id := uuid.New()
	return encryptors.AESEncrypt(id.String()[:aes.BlockSize])
}

func validateID(id string) (bool, error) {
	data, err := encryptors.AESDecrypt(id)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return regexp.Match(`\w{8}-\w{4}-\w{2}`, data)
}
