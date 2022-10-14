package handlers

import (
	"go-url-shortener/internal/storage"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const route = "/api/user/urls"

func TestGetUserLinks(t *testing.T) {
	tests := []struct {
		name    string
		storage []storage.ShortURL
		want    httpRes
	}{
		{
			name: "No stored links",
			want: httpRes{
				code:        http.StatusNoContent,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "One stored link",
			storage: []storage.ShortURL{{
				ID:      "id",
				URL:     "url",
				UID:     UserID,
				Deleted: false,
			}},
			want: httpRes{
				code:        http.StatusOK,
				resp:        `[{"short_url":"http://localhost:8080/id","original_url":"url"}]`,
				contentType: "application/json",
			},
		},
	}

	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := storage.NewMemoryRepo()
			if _, err := r.Add(tt.storage); err != nil {
				t.Fatal(err)
			}

			ts := getTestServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, http.MethodGet, route, "")
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.resp, body)
		})
	}
}

func TestDeleteUserLinks(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    httpRes
		wantErr string
	}{
		{
			name: "Missing body",
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        "you provided an incorrect IDs list format",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Malformed body",
			data: "test",
			want: httpRes{
				code:        http.StatusBadRequest,
				resp:        "you provided an incorrect IDs list format",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Correct body",
			data: `["id1","id2"]`,
			want: httpRes{code: http.StatusAccepted},
		},
	}

	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	ts := getTestServer(nil)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodDelete, route, tt.data)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.resp, body)
		})
	}
}
