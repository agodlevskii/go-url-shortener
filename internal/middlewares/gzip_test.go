package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	type (
		args struct {
			ct string
			cl string
		}
		want struct {
			ce     string
			reader string
		}
	)
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Missing content-type header",
			args: args{cl: "1500"},
			want: want{reader: "*httptest.ResponseRecorder"},
		},
		{
			name: "Incorrect size header",
			args: args{ct: "gzip"},
			want: want{reader: "*httptest.ResponseRecorder"},
		},
		{
			name: "Correct headers",
			args: args{
				ct: "gzip",
				cl: "1500",
			},
			want: want{
				ce:     "gzip",
				reader: "respwriters.GzipWriter",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ce := w.Header().Get("Content-Encoding")
				assert.Equal(t, tt.want.ce, ce)
				assert.Equal(t, tt.want.reader, reflect.TypeOf(w).String())
			})

			req := httptest.NewRequest(http.MethodGet, BaseURL, nil)
			req.Header.Add("Accept-Encoding", tt.args.ct)
			req.Header.Add("Content-Length", tt.args.cl)

			handler := Compress(next)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

func TestDecompress(t *testing.T) {
	type want struct {
		reader string
		writer string
	}

	tests := []struct {
		name string
		ct   string
		body string
		want want
	}{
		{
			name: "Missing content-type header",
			want: want{writer: "*httptest.ResponseRecorder"},
		},
		{
			name: "Correct content-type header",
			ct:   "gzip",
			want: want{
				reader: "*gzip.Reader",
				writer: "*httptest.ResponseRecorder",
			},
			body: "req_body",
		},
		{
			name: "Incorrect content-type header",
			ct:   "zip",
			want: want{writer: "*httptest.ResponseRecorder"},
			body: "req_body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.want.writer, reflect.TypeOf(w).String())
				if tt.want.reader != "" {
					assert.Equal(t, tt.want.reader, reflect.TypeOf(r.Body).String())
				}
			})

			req := httptest.NewRequest(http.MethodPost, BaseURL, io.NopCloser(bytes.NewBufferString(tt.body)))
			req.Header.Add("Content-Encoding", tt.ct)

			handler := Decompress(next)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
