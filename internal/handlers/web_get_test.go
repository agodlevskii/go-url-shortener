package handlers

import (
	"github.com/stretchr/testify/assert"
	"go-url-shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetFullURL(t *testing.T) {
	type (
		want struct {
			code        int
			resp        string
			contentType string
			location    string
		}
		testCase struct {
			name    string
			want    want
			id      string
			storage []storage.ShortURL
		}
	)

	tests := []testCase{
		{
			name: "No ID query parameter",
			id:   "",
			want: want{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Incorrect ID parameter value",
			id:   "foo",
			want: want{
				code:        http.StatusBadRequest,
				resp:        http.StatusText(http.StatusBadRequest),
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Correct ID parameter value",
			id:      "googl",
			storage: []storage.ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}},
			want: want{
				code:        http.StatusTemporaryRedirect,
				resp:        `https://google.com`,
				contentType: "text/html; charset=utf-8",
				location:    "https://google.com",
			},
		},
	}

	if err := os.Chdir("../../"); err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := storage.NewMemoryRepo()
			if _, err := db.Add(tt.storage); err != nil {
				t.Error(err)
			}

			r := NewShortenerRouter(db)
			ts := httptest.NewServer(r)
			defer ts.Close()

			path := "/"
			if tt.id != "" {
				path = path + tt.id
			}

			resp, body := testGetRequest(t, ts, path)
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))

			if tt.want.resp != "" {
				assert.Contains(t, body, tt.want.resp)
			}

			defer resp.Body.Close()
		})
	}
}
