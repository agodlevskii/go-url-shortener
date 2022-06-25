package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var UserID = "7190e4d4-fd9c-4b"

func TestMemoRepo_Add(t *testing.T) {
	type fields struct {
		db map[string]UrlRes
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
			fields:  fields{db: map[string]UrlRes{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}

			err := m.Add(UserID, tt.args.id, tt.args.url)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestMemoRepo_Clear(t *testing.T) {
	type fields struct {
		db map[string]UrlRes
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Correct clean",
			fields: fields{db: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			}},
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
		db map[string]UrlRes
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
			name: "Missing ID",
			fields: fields{db: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			}},
			args:    args{id: "foo"},
			wantErr: true,
		},
		{
			name: "Existing ID",
			fields: fields{db: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			}},
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

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, url)
		})
	}
}

func TestMemoRepo_Has(t *testing.T) {
	type fields struct {
		db map[string]UrlRes
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Missing ID",
			fields: fields{db: map[string]UrlRes{
				"bar": {
					url: "https://google.com",
					uid: UserID,
				},
			}},
			args: args{id: "foo"},
			want: false,
		},
		{
			name: "Existing ID",
			fields: fields{db: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			}},
			args: args{id: "googl"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemoRepo{
				db: tt.fields.db,
			}

			has, err := m.Has(tt.args.id)

			assert.Equal(t, tt.want, has)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

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
	type fields struct {
		filename string
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
			fields:  fields{filename: "testfile"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileRepo{
				filename: tt.fields.filename,
			}

			err := f.Add(UserID, tt.args.id, tt.args.url)
			assert.Equal(t, tt.wantErr, err != nil)

			f.Clear()
		})
	}
}

func TestFileRepo_Clear(t *testing.T) {
	type fields struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Successful clean",
			fields: fields{filename: "testfile"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileRepo{
				filename: tt.fields.filename,
			}

			f.Add(UserID, "googl", "https://google.com")
			has, err := f.Has("googl")
			assert.Equal(t, true, has)
			assert.Equal(t, false, err != nil)

			f.Clear()
			has, err = f.Has("googl")
			assert.Equal(t, false, has)
			assert.Equal(t, false, err != nil)
		})
	}
}

func TestFileRepo_Get(t *testing.T) {
	type fields struct {
		filename string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		data    map[string]UrlRes
		wantErr bool
	}{
		{
			name:   "Missing ID",
			fields: fields{filename: "testfile"},
			data: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			},
			args:    args{id: "foo"},
			wantErr: true,
		},
		{
			name:   "Existing ID",
			fields: fields{filename: "testfile"},
			data: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			},
			args:    args{id: "googl"},
			want:    "https://google.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileRepo{
				filename: tt.fields.filename,
			}

			for id, res := range tt.data {
				f.Add(UserID, id, res.url)
			}

			got, err := f.Get(tt.args.id)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.id)

			if len(tt.data) > 0 {
				f.Clear()
			}
		})
	}
}

func TestFileRepo_Has(t *testing.T) {
	type fields struct {
		filename string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		data    map[string]UrlRes
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "Missing ID",
			fields: fields{filename: "testfile"},
			args:   args{id: "foo"},
			want:   false,
		},
		{
			name:   "Existing ID",
			fields: fields{filename: "testfile"},
			data: map[string]UrlRes{
				"googl": {
					url: "https://google.com",
					uid: UserID,
				},
			},
			args: args{id: "googl"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileRepo{
				filename: tt.fields.filename,
			}

			for id, res := range tt.data {
				f.Add(UserID, id, res.url)
			}

			has, err := f.Has(tt.args.id)
			assert.Equal(t, tt.want, has)
			assert.Equal(t, tt.wantErr, err != nil)

			if len(tt.data) > 0 {
				f.Clear()
			}
		})
	}
}
