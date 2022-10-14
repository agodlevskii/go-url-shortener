package generators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateString(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Correct size",
			args: args{size: 3},
			want: 3,
		},
		{
			name:    "Missing size",
			args:    args{size: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateString(tt.args.size)
			assert.Equal(t, tt.want, len(got))
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
