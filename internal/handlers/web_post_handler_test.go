package handlers

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go-url-shortener/configs"
	"go-url-shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_generateID(t *testing.T) {
	type args struct {
		db   storage.MemoRepo
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Defined size",
			args: args{db: storage.NewMemoryRepo(), size: 3},
			want: 3,
		},
		{
			name: "Undefined size",
			args: args{db: storage.NewMemoryRepo(), size: 0},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := generateID(tt.args.db, tt.args.size)
			got := len(res)
			assert.Equalf(t, tt.want, got, "generateID(%v)", tt.args.db, tt.args.size)
		})
	}
}

func TestShortenURL(t *testing.T) {
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

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Error(err)
				}
			}(resp.Body)
		})
	}
}
