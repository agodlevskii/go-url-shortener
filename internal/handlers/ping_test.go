package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-url-shortener/internal/storage"
	"net/http"
	"os"
	"testing"
)

type mockDb struct {
	mock.Mock
}

func (m *mockDb) Ping() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockDb) Add([]storage.ShortURL) ([]storage.ShortURL, error) {
	return nil, nil
}

func (m *mockDb) Clear() {}

func (m *mockDb) Get(string) (storage.ShortURL, error) {
	return storage.ShortURL{}, nil
}

func (m *mockDb) GetAll(string) ([]storage.ShortURL, error) {
	return nil, nil
}

func (m *mockDb) Has(id string) (bool, error) {
	return true, nil
}

func (m *mockDb) Delete(batch []storage.ShortURL) error {
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
			db := new(mockDb)
			db.On("Ping").Return(tt.resp)

			ts := getTestServer(db)
			defer ts.Close()

			resp, body := testRequest(t, ts, http.MethodGet, "/ping", "")
			defer resp.Body.Close()

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.resp, body)
		})
	}
}
