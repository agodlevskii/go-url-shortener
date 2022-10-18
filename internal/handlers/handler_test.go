package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go-url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type httpRes struct {
	code        int
	resp        string
	contentType string
	location    string
}

const (
	BaseURL        = "http://localhost:8080"
	UserIDEnc      = "4b529d6712a1d59f62a87dc4fa54f332"
	UserID         = "7190e4d4-fd9c-4b"
	UserCookieName = "user_id"
)

func TestNewShortenerRouter(t *testing.T) {
	ts := getTestServer(nil)
	defer ts.Close()

	resp, _ := testRequest(t, ts, http.MethodPut, "/", "")
	assert.Error(t, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, data string) (*http.Response, string) {
	rawURL := ts.URL + path
	purl, _ := url.Parse(rawURL)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	jar.SetCookies(purl, []*http.Cookie{
		{Name: UserCookieName, Value: UserIDEnc, Path: "/"},
	})

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	var body io.Reader
	if data != "" {
		body = io.NopCloser(bytes.NewBufferString(data))
	}

	req, err := http.NewRequest(method, rawURL, body)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func(Body io.ReadCloser) {
		if cErr := Body.Close(); cErr != nil {
			t.Error(cErr)
		}
	}(resp.Body)

	return resp, strings.TrimSpace(string(respBody))
}

func getTestServer(repo storage.Storager) *httptest.Server {
	if repo == nil {
		repo = storage.NewMemoryRepo()
	}
	r := NewShortenerRouter(repo)
	return httptest.NewServer(r)
}
