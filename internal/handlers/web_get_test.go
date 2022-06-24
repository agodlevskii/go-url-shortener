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
			storage map[string]map[string]string
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
				resp:        "no matching URL found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct ID parameter value",
			id:   "googl",
			storage: map[string]map[string]string{
				UserId: {"googl": "https://google.com"},
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				resp:        `https://google.com`,
				contentType: "text/html; charset=utf-8",
				location:    "https://google.com",
			},
		},
	}

	db := storage.NewMemoryRepo()
	r := NewShortenerRouter(db, "https://test.url")
	ts := httptest.NewServer(r)
	defer ts.Close()

	err := os.Chdir("../../")
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.storage) > 0 {
				for uid, urls := range tt.storage {
					for id, url := range urls {
						err := db.Add(uid, id, url)
						if err != nil {
							t.Error(err)
						}
					}
				}
			}

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

			if len(tt.storage) > 0 {
				db.Clear()
			}
		})
	}
}
