package handlers

import (
	"github.com/stretchr/testify/assert"
	"go-url-shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebPostHandler(t *testing.T) {
	type (
		want struct {
			code        int
			resp        string
			contentType string
		}
		testCase struct {
			name         string
			data         string
			checkInclude bool
			want         want
		}
	)

	args := struct {
		repo    *storage.MemoRepo
		baseURL string
	}{
		repo:    storage.NewMemoryRepo(),
		baseURL: "https://test.url",
	}

	tests := []testCase{
		{
			name: "Missing body",
			want: want{
				code:        http.StatusBadRequest,
				resp:        "The original URL is missing. Please attach it to the request body.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty body",
			want: want{
				code:        http.StatusBadRequest,
				resp:        "The original URL is missing. Please attach it to the request body.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:         "Correct body",
			data:         "https://google.com",
			checkInclude: true,
			want: want{
				code:        http.StatusCreated,
				resp:        args.baseURL,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	r := NewShortenerRouter(args.repo, args.baseURL)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testPostRequest(t, ts, "/", tt.data)
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			if tt.checkInclude {
				assert.Contains(t, body, tt.want.resp)
			} else {
				assert.Equal(t, tt.want.resp, body)
			}

			defer resp.Body.Close()
		})
	}
}

func TestAPIPostHandler(t *testing.T) {
	type (
		want struct {
			code        int
			resp        string
			contentType string
		}
	)

	args := struct {
		repo    storage.Storager
		baseURL string
	}{
		repo:    storage.NewMemoryRepo(),
		baseURL: "https://test.url",
	}

	tests := []struct {
		name         string
		want         want
		data         string
		checkInclude bool
	}{
		{
			name: "Missing body",
			want: want{
				code:        http.StatusBadRequest,
				resp:        "You provided an incorrect URL request.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty body",
			want: want{
				code:        http.StatusBadRequest,
				resp:        "You provided an incorrect URL request.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:         "Correct body",
			data:         `{ "url": "https://google.com" }`,
			checkInclude: true,
			want: want{
				code:        http.StatusCreated,
				resp:        `{"result":` + `"` + args.baseURL,
				contentType: "application/json",
			},
		},
	}

	r := NewShortenerRouter(args.repo, args.baseURL)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testPostRequest(t, ts, "/api/shorten", tt.data)
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			if tt.checkInclude {
				assert.Contains(t, body, tt.want.resp)
			} else {
				assert.Equal(t, tt.want.resp, body)
			}

			defer resp.Body.Close()
		})
	}
}
