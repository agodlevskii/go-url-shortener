// Package respwriters includes the custom response writers, e.g. compressed data writer.
package respwriters

import (
	"io"
	"net/http"
)

// GzipWriter provides an implementation of the http.ResponseWriter interface for the compressed data.
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write is required to implement the io.Writer interface.
func (gw GzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}
