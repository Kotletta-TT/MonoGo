package common

import (
	"bytes"
	"compress/gzip"
)

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
