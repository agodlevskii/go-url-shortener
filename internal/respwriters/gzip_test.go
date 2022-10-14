package respwriters

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGzipWriter_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Correct compression",
			args: args{b: []byte("test")},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
			if err != nil {
				t.Fatal(err)
			}
			defer gz.Close()

			gw := GzipWriter{Writer: gz}
			got, err := gw.Write(tt.args.b)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
