package handlers

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-url-shortener/internal/storage"
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	return resp, strings.TrimSpace(string(respBody))
}

func testGetRequest(t *testing.T, ts *httptest.Server, path string) (*http.Response, string) {
	return testRequest(t, ts, http.MethodGet, path, "")
}

func testPostRequest(t *testing.T, ts *httptest.Server, path, data string) (*http.Response, string) {
	return testRequest(t, ts, http.MethodPost, path, data)
}

func TestNewShortenerRouter(t *testing.T) {
	r := NewShortenerRouter(storage.NewMemoryRepo(), "https://test.url")
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, http.MethodPut, "/", "")
	assert.Error(t, errors.New("This HTTP method is not allowed."))
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	defer resp.Body.Close()
}
