package generators

import (
	"github.com/stretchr/testify/assert"
	"go-url-shortener/internal/storage"
	"testing"
)

func TestGenerateID(t *testing.T) {
	type args struct {
		db   storage.MemoRepo
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Defined size",
			args: args{db: storage.NewMemoryRepo(), size: 3},
			want: 3,
		},
		{
			name: "Undefined size",
			args: args{db: storage.NewMemoryRepo(), size: 0},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GenerateID(tt.args.db, tt.args.size)
			got := len(res)
			assert.Equalf(t, tt.want, got, "generateID(%v)", tt.args.db, tt.args.size)
			assert.Equal(t, tt.wantErr, err == nil)
		})
	}
}
