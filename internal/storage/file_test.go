package storage

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
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
			r, err := NewFileRepo(tt.args.filename)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equalf(t, tt.want, r.filename, "NewFileRepo(%v)", tt.args.filename)
			r.Clear()
		})
	}
}

func TestFileRepo_Add(t *testing.T) {
	t.Parallel()

	fname := "testfile_add"
	for _, tt := range getAddTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewFileRepo(fname)
			if err != nil {
				t.Fatal(err)
			}

			got, err := r.Add(tt.batch)
			assert.Equal(t, tt.want.id, got[0].ID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
	cleanWorkspace(fname)
}

func TestFileRepo_Get(t *testing.T) {
	t.Parallel()

	fname := "testfile_get"
	for _, tt := range getGetTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewFileRepo(fname)
			if err != nil {
				t.Fatal(err)
			}

			_, err = r.Add([]ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}})
			if err != nil {
				t.Fatal(err)
			}

			sURL, err := r.Get(tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, sURL.URL)
		})
	}
	cleanWorkspace(fname)
}

func TestFileRepo_Has(t *testing.T) {
	t.Parallel()

	fname := "testfile_has"
	for _, tt := range getHasTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewFileRepo(fname)
			if err != nil {
				t.Fatal(err)
			}

			_, err = r.Add([]ShortURL{{ID: "googl", URL: "https://google.com", UID: UserID}})
			if err != nil {
				t.Fatal(err)
			}

			has, err := r.Has(tt.id)
			assert.Equal(t, tt.want, has)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
	cleanWorkspace(fname)
}

func TestFileRepo_Delete(t *testing.T) {
	t.Parallel()

	fname := "testfile_delete"
	for _, tt := range getDeleteTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewFileRepo(fname)
			if err != nil {
				t.Fatal(err)
			}

			_, err = r.Add(tt.batch)
			if err != nil {
				t.Fatal(err)
			}

			err = r.Delete(tt.batch)
			assert.Equal(t, tt.wantErr, err != nil)

			for _, sURL := range tt.batch {
				stored, err := r.Get(sURL.ID)
				if err != nil {
					r.Clear()
					t.Fatal(err)
				}

				assert.Equal(t, tt.wantDelState, stored.Deleted)
			}
		})
	}
	cleanWorkspace(fname)
}

func cleanWorkspace(fname string) {
	if err := os.Remove(fname); err != nil {
		log.Error(err)
	}
}
