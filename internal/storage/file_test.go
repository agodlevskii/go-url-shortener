package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFileRepo(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Filename is missing",
			args:    args{filename: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Filename is presented",
			args:    args{filename: "testfile"},
			want:    "testfile",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileRepo(tt.args.filename)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equalf(t, tt.want, got.filename, "NewFileRepo(%v)", tt.args.filename)
			got.Clear()
		})
	}
}

func TestFileRepo_Add(t *testing.T) {
	for _, tt := range getAddTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, err := NewFileRepo("testfile_add")
			if err != nil {
				t.Fatal(err)
			}

			got, err := repo.Add(tt.batch)
			assert.Equal(t, tt.want.id, got[0].ID)
			assert.Equal(t, tt.wantErr, err != nil)
			repo.Clear()
		})
	}
}

func TestFileRepo_Get(t *testing.T) {
	for _, tt := range getGetTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, err := NewFileRepo("testfile_get")
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Add([]ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}})
			if err != nil {
				t.Fatal(err)
			}
			sURL, err := repo.Get(tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, sURL.URL)
			repo.Clear()
		})
	}
}

func TestFileRepo_Has(t *testing.T) {
	for _, tt := range getHasTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, err := NewFileRepo("testfile_has")
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Add([]ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}})
			if err != nil {
				t.Fatal(err)
			}

			has, err := repo.Has(tt.id)
			assert.Equal(t, tt.want, has)
			assert.Equal(t, tt.wantErr, err != nil)
			repo.Clear()
		})
	}
}
