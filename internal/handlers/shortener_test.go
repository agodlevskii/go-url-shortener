package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
)

func TestWebShortener(t *testing.T) {
	type testCase struct {
		name         string
		data         string
		checkInclude bool
		want         httpRes
	}

	tests := []testCase{
		{
			name: "Missing body",
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        apperrors.URLFormat,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty body",
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        apperrors.URLFormat,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:         "Correct body",
			data:         "https://google.com",
			checkInclude: true,
			want: httpRes{
				code:        http.StatusCreated,
				resp:        BaseURL,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	ts := getTestServer(nil)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodPost, "/", tt.data)
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			if tt.checkInclude {
				assert.Contains(t, body, tt.want.resp)
			} else {
				assert.Equal(t, tt.want.resp, body)
			}

			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAPIShortener(t *testing.T) {
	tests := []struct {
		name         string
		want         httpRes
		data         string
		cookie       *http.Cookie
		checkInclude bool
	}{
		{
			name: "Missing cookie",
			data: `{ "url": "https://google.com" }`,
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        apperrors.UserID,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Missing body",
			cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        apperrors.URLFormat,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Empty body",
			cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        apperrors.URLFormat,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:         "Correct body",
			cookie:       &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
			data:         `{ "url": "https://google.com" }`,
			checkInclude: true,
			want: httpRes{
				code:        http.StatusCreated,
				resp:        BaseURL,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBufferString(tt.data))
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			w := httptest.NewRecorder()

			APIShortener(storage.NewMemoryRepo(), mockConfig{})(w, req)
			res := w.Result()
			b, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			body := strings.Trim(string(b), "\n")
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			if tt.checkInclude {
				assert.Contains(t, body, tt.want.resp)
			} else {
				assert.Equal(t, tt.want.resp, body)
			}

			if err = res.Body.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWebGetFullURL(t *testing.T) {
	type testCase struct {
		name   string
		want   httpRes
		id     string
		stored []storage.ShortURL
	}

	tests := []testCase{
		{
			name: "No ID query parameter",
			id:   "",
			want: httpRes{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Incorrect ID parameter value",
			id:   "foo",
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        http.StatusText(http.StatusBadRequest),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Correct ID parameter value",
			id:     "google",
			stored: []storage.ShortURL{{ID: "google", URL: "https://google.com", UID: UserID}},
			want: httpRes{
				code:        http.StatusTemporaryRedirect,
				resp:        `https://google.com`,
				contentType: "text/html; charset=utf-8",
				location:    "https://google.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := storage.NewMemoryRepo()
			if _, err := db.Add(context.Background(), tt.stored); err != nil {
				t.Fatal(err)
			}

			ts := getTestServer(db)
			defer ts.Close()

			path := "/"
			if tt.id != "" {
				path += tt.id
			}

			resp, body := testRequest(t, ts, http.MethodGet, path, "")
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))

			if tt.want.resp != "" {
				assert.Contains(t, body, tt.want.resp)
			}

			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAPIBatchShortener(t *testing.T) {
	type args struct {
		cookie *http.Cookie
		body   []BatchReqData
	}
	type want struct {
		code int
		cp   string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Missing cookie",
			want: want{code: http.StatusBadRequest},
		},
		{
			name: "Incorrect cookie",
			args: args{cookie: &http.Cookie{Name: UserCookieName, Value: "bad_cookie", Path: "/"}},
			want: want{code: http.StatusBadRequest},
		},
		{
			name: "Incorrect body",
			args: args{cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"}},
			want: want{code: http.StatusBadRequest},
		},
		{
			name: "Correct body",
			args: args{
				cookie: &http.Cookie{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
				body: []BatchReqData{
					{CorrelationID: "google", OriginalURL: "https://google.com"},
					{CorrelationID: "facebook", OriginalURL: "https://facebook.com"},
				},
			},
			want: want{
				code: http.StatusCreated,
				cp:   "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(b))
			if tt.args.cookie != nil {
				req.AddCookie(tt.args.cookie)
			}
			w := httptest.NewRecorder()

			APIBatchShortener(storage.NewMemoryRepo(), mockConfig{})(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if err = res.Body.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
