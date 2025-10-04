package headers

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const (
	crlf = "\r\n"
)

var headerKeyRegex = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+$`)

var ErrMissingColon = errors.New("invalid header: missing colon")
var ErrSpaceBeforeColon = errors.New("invalid header: space before colon")
var ErrInvalidKey = errors.New("invalid header key: key contains invalid characters")

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
		return 0, false, ErrMissingColon
	}

	// make sure no spaces before colon
	if strings.ContainsAny(trimmed[:colonIdx], " \t") {
		return 0, false, ErrSpaceBeforeColon
	}

	key := trimmed[:colonIdx]
	if !isValidHeaderKey(key) {
		return 0, false, ErrInvalidKey
	}
	val := strings.TrimSpace(trimmed[colonIdx+1:])
	keyInLowerCase := strings.ToLower(key)

	existingValues, ok := h[keyInLowerCase]
	if ok {
		h[keyInLowerCase] = fmt.Sprintf("%s, %s", existingValues, val)
	} else {
		h[keyInLowerCase] = val
	}

	// consumed line + CRLF
	return idx + len(crlf), false, nil
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	return h[key]
}

func isValidHeaderKey(headerKey string) bool {
	return headerKeyRegex.MatchString(headerKey)
}
