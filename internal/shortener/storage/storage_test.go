package storage

import (
	"testing"
)

func TestAddURLToStorage(t *testing.T) {
	type args struct {
		repo Storager
		id   string
		url  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Correct URL",
			args: args{
				repo: NewMemoryRepo(nil),
				id:   "googl",
				url:  "https://google.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddURLToStorage(tt.args.repo, tt.args.id, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("AddURLToStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetURLFromStorage(t *testing.T) {
	type args struct {
		repo Storager
		id   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Missing ID",
			args: args{
				repo: NewMemoryRepo(nil),
				id:   "googl",
			},
			wantErr: true,
		},
		{
			name: "Existing ID",
			args: args{
				repo: NewMemoryRepo(map[string]string{"googl": "https://google.com"}),
				id:   "googl",
			},
			want:    "https://google.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetURLFromStorage(tt.args.repo, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURLFromStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetURLFromStorage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoRepo_Add(t *testing.T) {
	type fields struct {
		db map[string]string
	}
	type args struct {
		id  string
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Correct URL",
			args: args{
				id:  "googl",
				url: "https://google.com",
			},
			fields:  fields{db: map[string]string{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}
			if err := m.Add(tt.args.id, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoRepo_Clear(t *testing.T) {
	type fields struct {
		db map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Correct clean",
			fields: fields{db: map[string]string{"googl": "https://google.com"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}
			m.Clear()

			if len(m.db) > 0 {
				t.Errorf("Clear(): expect the repo to be empty, but found #{lenm.db) elements")
			}
		})
	}
}

func TestMemoRepo_Get(t *testing.T) {
	type fields struct {
		db map[string]string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Missing ID",
			fields:  fields{db: map[string]string{"googl": "https://google.com"}},
			args:    args{id: "foo"},
			wantErr: true,
		},
		{
			name:    "Existing ID",
			fields:  fields{db: map[string]string{"googl": "https://google.com"}},
			args:    args{id: "googl"},
			want:    "https://google.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}
			got, err := m.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoRepo_Has(t *testing.T) {
	type fields struct {
		db map[string]string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Missing ID",
			fields: fields{db: map[string]string{"googl": "https://google.com"}},
			args:   args{id: "foo"},
			want:   false,
		},
		{
			name:   "Existing ID",
			fields: fields{db: map[string]string{"googl": "https://google.com"}},
			args:   args{id: "googl"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}
			if got := m.Has(tt.args.id); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoRepo_Remove(t *testing.T) {
	type fields struct {
		db map[string]string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Missing ID",
			fields:  fields{db: map[string]string{"googl": "https://google.com"}},
			args:    args{id: "foo"},
			wantErr: true,
		},
		{
			name:    "Existing ID",
			fields:  fields{db: map[string]string{"googl": "https://google.com"}},
			args:    args{id: "googl"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}
			if err := m.Remove(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
