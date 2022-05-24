package storage

import (
	"testing"
)

func TestAddUrlToStorage(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "URL without a prefix",
			args: args{url: "google.com"},
			want: "googl",
		},
		{
			name: "URL with a prefix",
			args: args{url: "https://google.com"},
			want: "googl",
		},
		{
			name: "Empty argument",
			args: args{url: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddUrlToStorage(tt.args.url); got != tt.want {
				t.Errorf("AddUrlToStorage() = %v, want %v", got, tt.want)
			}

			if Storage[tt.want] != tt.args.url && tt.want != "" {
				t.Errorf(`Expected storage value for "#{tt.args.url}" to be "#{tt.want}", but got ""`)
			}
		})
	}
}

func TestGetUrlFromStorage(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		storage map[string]string
	}{
		{
			name:    "Empty ID",
			args:    args{id: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Missing ID value",
			args:    args{id: "foo"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Existing ID value",
			args:    args{id: "googl"},
			want:    "https://google.com",
			wantErr: false,
			storage: map[string]string{"googl": "https://google.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.storage) > 0 {
				for k, v := range tt.storage {
					Storage[k] = v
				}
			}

			got, err := GetUrlFromStorage(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrlFromStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUrlFromStorage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
