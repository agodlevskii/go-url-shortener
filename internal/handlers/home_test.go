package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestGetHomePage(t *testing.T) {
	tests := []struct {
		name string
		want httpRes
	}{
		{
			name: "Homepage template",
			want: httpRes{
				code:        http.StatusOK,
				resp:        "<title>Go URL Shortener</title>",
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	ts := getTestServer(nil)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodGet, "/", "")
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Contains(t, body, tt.want.resp)
		})
	}
}
