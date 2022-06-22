package respwriters

import (
	"io"
	"net/http"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (gw GzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}
