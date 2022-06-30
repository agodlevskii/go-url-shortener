package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoRepo_Add(t *testing.T) {
	tests := getAddTestCases(NewMemoryRepo())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			err := r.Add(UserID, tt.args.id, tt.args.url)
			assert.Equal(t, tt.wantErr, err != nil)
			r.Clear()
		})
	}
}

func TestMemoRepo_AddAll(t *testing.T) {
	tests := getAddAllTestCases(NewMemoryRepo())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			err := r.AddAll(UserID, tt.batch)
			assert.Equal(t, tt.wantErr, err != nil)
			r.Clear()
		})
	}
}

func TestMemoRepo_Clear(t *testing.T) {
	repo := NewMemoryRepo()
	err := repo.Add("googl", "https://google.com", UserID)
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

func TestMemoRepo_Get(t *testing.T) {
	repo := NewMemoryRepo()
	err := repo.Add(UserID, "googl", "https://google.com")
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

func TestMemoRepo_Has(t *testing.T) {
	repo := NewMemoryRepo()
	err := repo.Add(UserID, "googl", "https://google.com")
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
