package utils

import (
	"testing"
)

func TestIsURLStringValid(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Incorrect URL",
			args: args{rawURL: "URL"},
			want: false,
		},
		{
			name: "URL without protocol",
			args: args{rawURL: "google.com"},
			want: false,
		},
		{
			name: "Correct URL",
			args: args{rawURL: "https://google.com"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURLStringValid(tt.args.rawURL); got != tt.want {
				t.Errorf("IsURLStringValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
