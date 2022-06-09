package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			assert.Equal(t, tt.wantErr, err != nil)
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

			assert.Equal(t, tt.wantErr, err != nil)
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

			err := f.Add(tt.args.id, tt.args.url)
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

			f.Add("googl", "https://google.com")
			assert.Equal(t, true, f.Has("googl"))

			f.Clear()
			assert.Equal(t, false, f.Has("googl"))
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
		data    map[string]string
		wantErr bool
	}{
		{
			name:    "Missing ID",
			fields:  fields{filename: "testfile"},
			data:    map[string]string{"googl": "https://google.com"},
			args:    args{id: "foo"},
			wantErr: true,
		},
		{
			name:    "Existing ID",
			fields:  fields{filename: "testfile"},
			data:    map[string]string{"googl": "https://google.com"},
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

			for k, v := range tt.data {
				f.Add(k, v)
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
		name   string
		fields fields
		data   map[string]string
		args   args
		want   bool
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
			data:   map[string]string{"googl": "https://google.com"},
			args:   args{id: "googl"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileRepo{
				filename: tt.fields.filename,
			}

			for k, v := range tt.data {
				f.Add(k, v)
			}

			assert.Equalf(t, tt.want, f.Has(tt.args.id), "Has(%v)", tt.args.id)

			if len(tt.data) > 0 {
				f.Clear()
			}
		})
	}
}
