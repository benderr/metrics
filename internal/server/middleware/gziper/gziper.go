package gziper

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"

	_ "github.com/go-chi/chi/middleware"
)

func New(level int, contentTypes ...string) *GzipCompressor {
	allowedTypes := make(map[string]struct{})
	for _, t := range contentTypes {
		allowedTypes[t] = struct{}{}
	}

	return &GzipCompressor{
		level:        level,
		allowedTypes: allowedTypes,
	}
}

type GzipCompressor struct {
	level        int
	allowedTypes map[string]struct{}
}

func (g *GzipCompressor) TransformWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cw := &compressWriter{
			ResponseWriter: w,
			w:              w,
			contentTypes:   g.allowedTypes,
			compressable:   false,
		}

		encoder, _ := g.getWriter(r.Header, w)
		if encoder != nil {
			cw.w = encoder
		}

		//defer cleanup()
		defer cw.Close()

		next.ServeHTTP(cw, r)
	})
}

func (g *GzipCompressor) TransformReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(w, r)
	})
}

func (g *GzipCompressor) getWriter(h http.Header, w io.Writer) (io.Writer, error) {
	header := h.Get("Accept-Encoding")

	accepted := strings.Split(strings.ToLower(header), ",")

	for _, v := range accepted {
		if v == "gzip" {
			wr, err := gzip.NewWriterLevel(w, g.level)
			if err != nil {
				return nil, err
			}
			return wr, nil
		}
	}

	return nil, nil
}

type compressWriter struct {
	http.ResponseWriter

	w            io.Writer
	contentTypes map[string]struct{}
	wroteHeader  bool
	compressable bool
}

func (cw *compressWriter) isCompressable() bool {
	contentType := cw.Header().Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = contentType[0:idx]
	}

	if _, ok := cw.contentTypes[contentType]; ok {
		return true
	}

	return false
}

func (cw *compressWriter) WriteHeader(code int) {
	if cw.wroteHeader {
		cw.ResponseWriter.WriteHeader(code)
		return
	}
	cw.wroteHeader = true
	defer cw.ResponseWriter.WriteHeader(code)

	if cw.Header().Get("Content-Encoding") != "" {
		return
	}

	if !cw.isCompressable() {
		cw.compressable = false
		return
	}

	cw.compressable = true
	cw.Header().Set("Content-Encoding", "gzip")
	cw.Header().Set("Vary", "Accept-Encoding")
	cw.Header().Del("Content-Length")
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	if !cw.wroteHeader {
		cw.WriteHeader(http.StatusOK)
	}

	return cw.writer().Write(p)
}

func (cw *compressWriter) Close() error {
	if c, ok := cw.writer().(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("chi/middleware: io.WriteCloser is unavailable on the writer")
}

func (cw *compressWriter) writer() io.Writer {
	if cw.compressable {
		return cw.w
	} else {
		return cw.ResponseWriter
	}
}

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		ReadCloser: r,
		zr:         zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.ReadCloser.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
