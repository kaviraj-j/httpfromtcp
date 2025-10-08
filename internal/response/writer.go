package response

import (
	"fmt"
	"httpfromtcp/kaviraj-j/internal/headers"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reasonPhrase := ""
	switch statusCode {
	case StatusOk:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	statusLineStr := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	_, err := w.writer.Write([]byte(statusLineStr))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		headerData := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(headerData))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}
