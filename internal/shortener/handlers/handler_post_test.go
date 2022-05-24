package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenerPostHandler(t *testing.T) {
	type want struct {
		code        int
		resp        string
		contentType string
	}
	type testCase struct {
		name string
		args handlerArgs
		want want
	}

	tests := []testCase{
		{
			name: "Missing body",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want: want{
				code:        400,
				resp:        "The original URL is missing. Please attach it to the request body.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty body",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(([]byte)(``))),
			},
			want: want{
				code:        400,
				resp:        "The original URL is missing. Please attach it to the request body.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct body",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(([]byte)(`https://google.com`))),
			},
			want: want{
				code:        201,
				resp:        "googl",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(ShortenerPostHandler)
			h.ServeHTTP(tt.args.w, tt.args.r)

			res := tt.args.w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf(`Expected status code is %d, but got %d`, tt.want.code, res.StatusCode)
			}

			ct := res.Header.Get("Content-Type")
			if ct != tt.want.contentType {
				t.Errorf(`Expected content type is "%s", but got "%s"`, tt.want.contentType, ct)
			}

			defer res.Body.Close()
			b, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			bs := strings.TrimSpace(string(b))
			if bs != tt.want.resp {
				t.Errorf(`Expected response is "%s", but got "%s"`, tt.want.resp, bs)
			}
		})
	}
}
