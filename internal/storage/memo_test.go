package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoRepo_Add(t *testing.T) {
	for _, tt := range getAddTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewMemoryRepo()
			got, err := r.Add(tt.batch)
			assert.Equal(t, tt.want.id, got[0].ID)
			assert.Equal(t, tt.wantErr, err != nil)
			r.Clear()
		})
	}
}

func TestMemoRepo_Clear(t *testing.T) {
	for _, tt := range getClearTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewMemoryRepo()
			_, err := r.Add([]ShortURL{{ID: "googl", URL: "https://google.com"}})
			if err != nil {
				t.Fatal(err)
			}
			r.Clear()

			res, err := r.GetAll(UserID)
			if err != nil {
				t.Error(err)
			}
			assert.Zero(t, len(res))
		})
	}
}

func TestMemoRepo_Get(t *testing.T) {
	for _, tt := range getGetTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewMemoryRepo()
			_, err := r.Add([]ShortURL{{ID: "googl", URL: "https://google.com"}})
			if err != nil {
				t.Fatal(err)
			}

			sURL, err := r.Get(tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, sURL.URL)
		})
	}
}

func TestMemoRepo_Has(t *testing.T) {
	for _, tt := range getHasTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewMemoryRepo()
			_, err := r.Add([]ShortURL{{ID: "googl", URL: "https://google.com"}})
			if err != nil {
				t.Fatal(err)
			}

			has, err := r.Has(tt.id)
			assert.Equal(t, tt.want, has)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestMemoRepo_Delete(t *testing.T) {
	for _, tt := range getDeleteTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewMemoryRepo()
			_, err := r.Add(tt.batch)
			if err != nil {
				t.Fatal(err)
			}

			err = r.Delete(tt.batch)
			assert.Equal(t, tt.wantErr, err != nil)

			for _, sURL := range tt.batch {
				stored, err := r.Get(sURL.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.wantDelState, stored.Deleted)
			}
		})
	}
}
