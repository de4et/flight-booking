package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
)

type GzipCompressor struct {
	level int
}

func NewGzipCompressor(level int) *GzipCompressor {
	return &GzipCompressor{
		level: level,
	}
}

func (c *GzipCompressor) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&b, c.level)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}
	slog.Debug("Compressed ", "was", len(data), "then", len(b.Bytes()))
	return b.Bytes(), nil
}

func (c *GzipCompressor) Decompress(compressed []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	res, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
