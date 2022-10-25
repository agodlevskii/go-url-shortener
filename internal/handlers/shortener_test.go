package handlers

import (
	"context"
	"net/http"
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
		checkInclude bool
	}{
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
			data:         `{ "url": "https://google.com" }`,
			checkInclude: true,
			want: httpRes{
				code:        http.StatusCreated,
				resp:        `{"result":` + `"` + BaseURL,
				contentType: "application/json",
			},
		},
	}

	ts := getTestServer(nil)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", tt.data)

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
