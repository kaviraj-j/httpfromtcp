package response

import (
	"fmt"
	"httpfromtcp/kaviraj-j/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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
	_, err := w.Write([]byte(statusLineStr))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := map[string]string{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}

	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for k, v := range h {
		headerData := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(headerData))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
