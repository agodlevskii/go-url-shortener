package storage

import (
	"github.com/stretchr/testify/assert"
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
				repo: NewMemoryRepo(),
				id:   "googl",
				url:  "https://google.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddURLToStorage(tt.args.repo, tt.args.id, tt.args.url)
			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err, tt.wantErr)
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
		storage map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "Missing ID",
			args: args{
				repo: NewMemoryRepo(),
				id:   "googl",
			},
			wantErr: true,
		},
		{
			name: "Existing ID",
			args: args{
				repo: NewMemoryRepo(),
				id:   "googl",
			},
			storage: map[string]string{"googl": "https://google.com"},
			want:    "https://google.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.storage != nil {
				for k, v := range tt.storage {
					tt.args.repo.Add(k, v)
				}
			}
			url, err := GetURLFromStorage(tt.args.repo, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, url)
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

			err := m.Add(tt.args.id, tt.args.url)
			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err, tt.wantErr)
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

			assert.Zero(t, len(m.db))
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
			url, err := m.Get(tt.args.id)

			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, url)
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

			assert.Equal(t, tt.want, m.Has(tt.args.id))
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

			err := m.Remove(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err, tt.wantErr)
			}
		})
	}
}
