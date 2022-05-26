package handlers

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path, data string) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var body io.Reader
	if data != "" {
		body = io.NopCloser(bytes.NewBufferString(data))
	}

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, strings.TrimSpace(string(respBody))
}

func testGetRequest(t *testing.T, ts *httptest.Server, path string) (*http.Response, string) {
	return testRequest(t, ts, http.MethodGet, path, "")
}

func testPostRequest(t *testing.T, ts *httptest.Server, path, data string) (*http.Response, string) {
	return testRequest(t, ts, http.MethodPost, path, data)
}

func TestNewShortenerRouter(t *testing.T) {
	r := NewShortenerRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, http.MethodPut, "/", "")
	assert.Error(t, errors.New("This HTTP method is not allowed."))
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestShortenUrl(t *testing.T) {
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
				code:        400,
				resp:        "The original URL is missing. Please attach it to the request body.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty body",
			want: want{
				code:        400,
				resp:        "The original URL is missing. Please attach it to the request body.",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:         "Correct body",
			data:         "https://google.com",
			checkInclude: true,
			want: want{
				code:        201,
				resp:        "XVlBzgb",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	r := NewShortenerRouter()
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
		})
	}
}

func TestShortenerGetHandler(t *testing.T) {
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
			storage map[string]string
		}
	)

	tests := []testCase{
		{
			name: "No ID query parameter",
			id:   "",
			want: want{
				code:        201,
				resp:        index,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Incorrect ID parameter value",
			id:   "foo",
			want: want{
				code:        400,
				resp:        "no matching URL found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "Correct ID parameter value",
			id:      "googl",
			storage: map[string]string{"googl": "https://google.com"},
			want: want{
				code:        307,
				resp:        "https://google.com",
				contentType: "text/plain; charset=utf-8",
				location:    "https://google.com",
			},
		},
	}

	r := NewShortenerRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.storage) > 0 {
				db.Init(tt.storage)
			}

			path := "/"
			if tt.id != "" {
				path = path + tt.id
			}

			resp, body := testGetRequest(t, ts, path)
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			assert.Equal(t, tt.want.resp, body)

			if len(tt.storage) > 0 {
				db.Clear()
			}
		})
	}
}
