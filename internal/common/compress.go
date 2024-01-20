package common

import (
	"bytes"
	"compress/gzip"
)

// GzipCompress compresses the given data using gzip algorithm.
//
// It takes a byte slice as input parameter and returns the compressed byte slice and an error if any.
func GzipCompress(data []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err = gzipWriter.Write(data); err != nil {
		return nil, err
	}
	if err = gzipWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
