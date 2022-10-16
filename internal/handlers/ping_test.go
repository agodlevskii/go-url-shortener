package handlers

import (
	"go-url-shortener/internal/storage"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) Ping() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockDB) Add([]storage.ShortURL) ([]storage.ShortURL, error) {
	return nil, nil
}

func (m *mockDB) Clear() {}

func (m *mockDB) Get(string) (storage.ShortURL, error) {
	return storage.ShortURL{}, nil
}

func (m *mockDB) GetAll(string) ([]storage.ShortURL, error) {
	return nil, nil
}

func (m *mockDB) Has(string) (bool, error) {
	return true, nil
}

func (m *mockDB) Delete([]storage.ShortURL) error {
	return nil
}

func (m *mockDB) Close() error {
	return nil
}

func TestPing(t *testing.T) {
	tests := []struct {
		name string
		resp bool
		want httpRes
	}{
		{
			name: "Ping success",
			resp: true,
			want: httpRes{
				code:        http.StatusOK,
				resp:        "DB is up and running",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Ping fail",
			want: httpRes{
				code:        http.StatusInternalServerError,
				resp:        http.StatusText(http.StatusInternalServerError),
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := new(mockDB)
			db.On("Ping").Return(tt.resp)

			ts := getTestServer(db)
			defer ts.Close()

			resp, body := testRequest(t, ts, http.MethodGet, "/ping", "")
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.resp, body)

			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
