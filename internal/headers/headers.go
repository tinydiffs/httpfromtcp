package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	headers := make(Headers)
	return headers
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	idx := bytes.Index(data, []byte(crlf))

	if idx == 0 {
		return 0, true, nil
	}
	if idx == -1 {
		return 0, false, nil
	}

	headerString := string(data[:idx])

	before, after, found := strings.Cut(headerString, ":")
	if !found {
		return 0, false, fmt.Errorf("format error: no colon found in header string")
	}
	
	field_name := before
	field_value := strings.TrimSpace(after)

	if strings.Contains(field_name, " ") {
		return 0, false, fmt.Errorf("format error: whitespace in field-name")
	}

	h[field_name] = field_value
	return len(headerString) + 2, false, nil
}