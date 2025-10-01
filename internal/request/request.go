package request

import (
	"bytes"
	"errors"
	"io"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	SEPARATOR = "\r\n"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

// errors
var ErrInvalidRequestData = errors.New("invalid request data")
var ErrInvalidRequestLine = errors.New("invalid request line")
var ErrReadingDataInDoneState = errors.New("reading request data in done state")
var ErrUnknownState = errors.New("unknown state")

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buffer := make([]byte, 1024)
	bufferLen := 0
	for !request.done() {
		n, err := reader.Read(buffer[bufferLen:])
		if err != nil {
			if errors.Is(err, io.EOF) && request.done() {
				return request, nil
			}
			return nil, err
		}
		bufferLen += n
		consumed, err := request.parse(buffer[:bufferLen])
		if err != nil {
			return nil, err
		}
		if consumed > 0 {
			copy(buffer, buffer[consumed:bufferLen])
			bufferLen -= consumed
		}
	}

	return request, nil
}

func (request *Request) parse(data []byte) (int, error) {
	switch request.state {
	case StateInit:
		// parse the line
		requestLine, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if consumed == 0 {
			return 0, nil
		}
		request.RequestLine = *requestLine
		request.state = StateDone
		return consumed, nil
	case StateDone:
		return 0, ErrReadingDataInDoneState
	default:
		return 0, ErrUnknownState
	}
}

func (request *Request) done() bool {
	return request.state == StateDone
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(SEPARATOR))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineData := data[:idx]
	read := idx + len(SEPARATOR)
	parts := bytes.Split(requestLineData, []byte(" "))
	if len(parts) < 3 {
		return nil, 0, ErrInvalidRequestLine
	}
	httpParts := bytes.Split(parts[2], []byte("/"))
	return &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}, read, nil
}
