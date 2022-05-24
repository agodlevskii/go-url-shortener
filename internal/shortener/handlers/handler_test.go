package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerArgs struct {
	w *httptest.ResponseRecorder
	r *http.Request
}

func TestShortenerHandler(t *testing.T) {
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
			name: "Non-supported method",
			args: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPut, "/", nil),
			},
			want: want{
				code:        405,
				resp:        "HTTP request method is not supported.",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(ShortenerHandler)
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
		})
	}
}
