package generator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/spec"
	"io"
	"strings"
)

type ByteVarWriter struct {
	io.Writer
	c int
}

func newByteVarWriter(w io.Writer) *ByteVarWriter {
	return &ByteVarWriter{Writer: w}
}

func (w *ByteVarWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	for n = range p {
		if w.c%18 == 0 {
			w.Writer.Write([]byte("\n"))
		}
		fmt.Fprintf(w.Writer, "0x%02x,", p[n])
		w.c++
	}

	n++

	return
}

func encodeSwaggerToByteVar(v *spec.Swagger) (string, error) {
	b := &strings.Builder{}
	w := newByteVarWriter(b)
	if err := encodeSwagger(v, w); err != nil {
		return "", err
	}
	return b.String(), nil
}

func encodeSwagger(v *spec.Swagger, w io.Writer) error {
	gzWriter := gzip.NewWriter(w)
	encoder := json.NewEncoder(gzWriter)
	err := encoder.Encode(v)
	if err = gzWriter.Flush(); err != nil {
		return err
	}
	return gzWriter.Close()
}

func gzDecode(data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)
	gzReader, err := gzip.NewReader(dataReader)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	return io.ReadAll(gzReader)
}
