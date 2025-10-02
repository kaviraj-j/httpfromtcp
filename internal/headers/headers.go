package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

const (
	crlf = "\r\n"
)

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// find CRLF
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		// not enough data yet
		return 0, false, nil
	}

	// if CRLF is right at the start, we reached end of headers
	if idx == 0 {
		return len(crlf), true, nil
	}

	// extract the line (without CRLF)
	line := string(data[:idx])

	// trim leading/trailing spaces from the whole line
	trimmed := strings.TrimSpace(line)

	// find colon
	colonIdx := strings.Index(trimmed, ":")
	if colonIdx == -1 {
		return 0, false, errors.New("invalid header: missing colon")
	}

	// make sure no spaces before colon
	if strings.ContainsAny(trimmed[:colonIdx], " \t") {
		return 0, false, errors.New("invalid header: space before colon")
	}

	key := trimmed[:colonIdx]
	val := strings.TrimSpace(trimmed[colonIdx+1:])

	h[key] = val

	// consumed line + CRLF
	return idx + len(crlf), false, nil
}
