package handlers

import (
	"go-url-shortener/internal/shortener/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenerGetHandler(t *testing.T) {
	type want struct {
		code        int
		resp        string
		contentType string
		location    string
	}
	type testCase struct {
		name    string
		args    handlerArgs
		want    want
		storage map[string]string
	}

	tests := []testCase{
		{
			name: "No ID query parameter",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: want{
				code:        201,
				resp:        index,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Empty ID parameter",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/?id=", nil),
			},
			want: want{
				code:        201,
				resp:        index,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Incorrect ID parameter value",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/?id=foo", nil),
			},
			want: want{
				code:        400,
				resp:        "the URL with associated ID is not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Correct ID parameter value",
			storage: map[string]string{"googl": "https://google.com"},
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/?id=googl", nil),
			},
			want: want{
				code:        307,
				resp:        "https://google.com",
				contentType: "text/plain; charset=utf-8",
				location:    "https://google.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.storage) > 0 {
				for k, v := range tt.storage {
					storage.Storage[k] = v
				}
			}

			h := http.HandlerFunc(ShortenerGetHandler)
			h.ServeHTTP(tt.args.w, tt.args.r)

			res := tt.args.w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf(`Expected status code is %d, but got %d`, tt.want.code, res.StatusCode)
			}

			ct := res.Header.Get("Content-Type")
			if ct != tt.want.contentType {
				t.Errorf(`Expected content type is "%s", but got "%s"`, tt.want.contentType, ct)
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			bs := strings.TrimSpace(string(b))
			if bs != tt.want.resp {
				t.Errorf(`Expected response is "%s", but got "%s"`, tt.want.resp, bs)
			}

			loc := res.Header.Get("Location")
			if tt.want.location != "" && loc != tt.want.location {
				t.Errorf(`Expected location is "%s", but got "%s"`, tt.want.location, loc)
			}

			if len(tt.storage) > 0 {
				for k := range tt.storage {
					delete(storage.Storage, k)
				}
			}
		})
	}
}
