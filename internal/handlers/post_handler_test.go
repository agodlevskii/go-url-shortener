package handlers

import (
	"github.com/stretchr/testify/assert"
	"go-url-shortener/configs"
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
				resp:        "http://" + configs.Host + ":" + configs.Port,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	r := NewShortenerRouter(storage.NewMemoryRepo())
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
		args struct {
			db storage.MemoRepo
		}
		want struct {
			code        int
			resp        string
			contentType string
		}
	)

	tests := []struct {
		name         string
		args         args
		want         want
		data         string
		checkInclude bool
	}{
		{
			name: "Missing body",
			args: args{db: storage.NewMemoryRepo()},
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
				resp:        `{"result":` + `"http://` + configs.Host + ":" + configs.Port,
				contentType: "application/json",
			},
		},
	}

	r := NewShortenerRouter(storage.NewMemoryRepo())
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
