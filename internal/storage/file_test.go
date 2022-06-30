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
		})
	}
}

func TestFileRepo_Add(t *testing.T) {
	repo, err := NewFileRepo("testfile")
	if err != nil {
		t.Fatal(err)
	}

	tests := getAddTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			err = r.Add(UserID, tt.args.id, tt.args.url)
			assert.Equal(t, tt.wantErr, err != nil)
			r.Clear()
		})
	}
}

func TestFileRepo_AddAll(t *testing.T) {
	repo, err := NewFileRepo("testfile")
	if err != nil {
		t.Fatal(err)
	}

	tests := getAddAllTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			err = r.AddAll(UserID, tt.batch)
			assert.Equal(t, tt.wantErr, err != nil)
			r.Clear()
		})
	}
}

func TestFileRepo_Clear(t *testing.T) {
	repo, err := NewFileRepo("testfile")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Add("googl", "https://google.com", UserID)
	if err != nil {
		t.Fatal(err)
	}

	tests := getClearTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			r.Clear()

			res, err := r.GetAll(UserID)
			if err != nil {
				t.Error(err)
			}
			assert.Zero(t, len(res))
		})
	}
}

func TestFileRepo_Get(t *testing.T) {
	repo, err := NewFileRepo("testfile")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Add(UserID, "googl", "https://google.com")
	if err != nil {
		t.Fatal(err)
	}

	tests := getGetTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			url, err := r.Get(tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, url)
		})
	}

	repo.Clear()
}

func TestFileRepo_Has(t *testing.T) {
	repo, err := NewFileRepo("testfile")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Add(UserID, "googl", "https://google.com")
	if err != nil {
		t.Fatal(err)
	}

	tests := getHasTestCases(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			has, err := r.Has(tt.id)
			assert.Equal(t, tt.want, has)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

	repo.Clear()
}
