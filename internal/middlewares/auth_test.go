package middlewares

import (
	"go-url-shortener/internal/encryptors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BaseURL        = "http://localhost:8080"
	UserIDEnc      = "4b529d6712a1d59f62a87dc4fa54f332"
	UserID         = "7190e4d4-fd9c-4b"
	UserCookieName = "user_id"
)

func TestAuthorize(t *testing.T) {
	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(UserCookieName)
		if err != nil {
			t.Error(err)
		}

		dec, err := encryptors.AESDecrypt(cookie.Value)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 16, len(dec))
	})

	req := httptest.NewRequest(http.MethodGet, BaseURL, nil)
	handler := Authorize(next)
	handler.ServeHTTP(httptest.NewRecorder(), req)
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name    string
		cookie  *http.Cookie
		want    string
		wantErr bool
	}{
		{
			name:   "Existing cookie",
			cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
			want:   UserID,
		},
		{
			name:    "Missing",
			wantErr: true,
		},
	}

	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.cookie != nil {
					http.SetCookie(w, tt.cookie)
				}

				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(200)
			}))
			defer ts.Close()

			do, err := ts.Client().Do(httptest.NewRequest(http.MethodGet, BaseURL, nil))
			if err != nil {
				return
			}
			defer do.Body.Close()

			got, err := GetUserID(do.Request)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
