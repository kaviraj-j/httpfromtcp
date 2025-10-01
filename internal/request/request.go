package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	SEPARATOR = "\r\n"
)

// errors
var ErrInvalidRequestData = errors.New("invalid request data")
var ErrInvalidRequestLine = errors.New("invalid request line")

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	data, err := io.ReadAll(reader)
	if err != nil {
		return &request, nil
	}
	dataStr := string(data)
	requestLine, dataStr, err := parseRequestLine(dataStr)
	if err != nil {
		return &request, err
	}
	request.RequestLine = *requestLine

	return &request, nil
}

func parseRequestLine(data string) (*RequestLine, string, error) {
	idx := strings.Index(data, SEPARATOR)
	if idx == -1 {
		return nil, data, ErrInvalidRequestData
	}
	requestLineStr := data[:idx]
	remainingData := data[idx+len(SEPARATOR):]
	parts := strings.Split(requestLineStr, " ")
	if len(parts) < 3 {
		return nil, remainingData, ErrInvalidRequestLine
	}

	return &RequestLine{
		HttpVersion:   parts[2][5:],
		Method:        parts[0],
		RequestTarget: parts[1],
	}, remainingData, nil
}
