package middlewares

import (
	"crypto/aes"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/configs"
	"go-url-shortener/internal/encryptors"
	"io"
	"net/http"
	"regexp"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(configs.UserCookieName)

		if err == nil {
			valid, err := validateId(cookie.Value)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}

			if valid {
				next.ServeHTTP(w, r)
				return
			}
		}

		newId, err := generateId()
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		cookie = &http.Cookie{Name: configs.UserCookieName, Value: newId, Path: "/"}
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
		next.ServeHTTP(w, r)
	})
}

func GetUserId(r *http.Request) (string, error) {
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

func generateId() (string, error) {
	id := uuid.New()
	return encryptors.AESEncrypt(id.String()[:aes.BlockSize])
}

func validateId(id string) (bool, error) {
	data, err := encryptors.AESDecrypt(id)
	if err != nil {
		log.Error(err)
		return false, nil
	}
	return regexp.Match(`\w{8}-\w{4}-\w{2}`, data)
}
