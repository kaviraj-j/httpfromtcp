package server

import (
	"httpfromtcp/kaviraj-j/internal/request"
	"httpfromtcp/kaviraj-j/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func (hErr *HandlerError) Write(w io.Writer) error {
	response.WriteStatusLine(w, hErr.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(hErr.Message)))
	w.Write([]byte(hErr.Message))
	return nil
}
