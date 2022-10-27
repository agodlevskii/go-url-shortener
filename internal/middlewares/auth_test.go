package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-url-shortener/internal/encryptors"
)

type mockConfig struct{}

func (m mockConfig) GetUserCookieName() string {
	return UserCookieName
}

const (
	BaseURL        = "http://localhost:8080"
	UserIDEnc      = "4b529d6712a1d59f62a87dc4fa54f332"
	UserID         = "7190e4d4-fd9c-4b"
	UserCookieName = "user_id"
)

func TestAuthorize(t *testing.T) {
	tests := []struct {
		name      string
		cookie    *http.Cookie
		want      string
		wantOther bool
	}{
		{
			name:      "Missing cookie",
			wantOther: true,
		},
		{
			name:   "Valid cookie",
			cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
			want:   UserID,
		},
		{
			name:      "Invalid cookie",
			cookie:    &http.Cookie{Name: UserCookieName, Value: "bad_cookie", Path: "/"},
			wantOther: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cookie, err := r.Cookie(UserCookieName)
				if err != nil {
					t.Fatal(err)
				}

				dec, err := encryptors.AESDecrypt(cookie.Value)
				if err != nil {
					t.Fatal(err)
				}

				if cookie.Value != UserIDEnc {
					assert.True(t, tt.wantOther)
				} else {
					assert.Equal(t, tt.want, string(dec))
				}

			})

			req := httptest.NewRequest(http.MethodGet, BaseURL, nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			handler := Authorize(mockConfig{})(next)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name    string
		cookie  *http.Cookie
		want    string
		wantErr bool
	}{
		{
			name:    "Missing cookie",
			wantErr: true,
		},
		{
			name:   "Valid cookie",
			cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
			want:   UserID,
		},
		{
			name:    "Invalid cookie",
			cookie:  &http.Cookie{Name: UserCookieName, Value: "bad_cookie", Path: "/"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, BaseURL, nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			got, err := GetUserID(mockConfig{}, req)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
