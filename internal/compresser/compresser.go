package compresser

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	decoder, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: decoder,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Add("Content-Encoding", "gzip")
	}

	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Write(b []byte) (n int, err error) {
	return cw.zw.Write(b)
}

func (cw *compressWriter) Close() error {
	return cw.zw.Close()
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		contentEncoding := r.Header.Get("Content-Encoding")
		needDecompress := strings.Contains(contentEncoding, "gzip")
		if needDecompress {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				log.Error().Err(err).Msgf("")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		if supportGzip {
			fmt.Println("Compress")
			cw := newCompressWriter(w)

			ow = cw
			defer cw.zw.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
