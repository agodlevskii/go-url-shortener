package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHomePage(t *testing.T) {
	tests := []struct {
		name string
		want httpRes
	}{
		{
			name: "Homepage message",
			want: httpRes{
				code:        http.StatusOK,
				resp:        "The URL shortener is up and running.",
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	ts := getTestServer(nil)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, http.MethodGet, "/", "")
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, body, tt.want.resp)

			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
