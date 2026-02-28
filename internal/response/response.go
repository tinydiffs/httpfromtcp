package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/tinydiffs/httpfromtcp/internal/headers"
)


type StatusCode int

const(
	Ok StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	reasonPhrase := ""

	switch statusCode {
	case Ok:
		reasonPhrase = "HTTP/1.1 200 OK" + crlf
	case BadRequest:
		reasonPhrase = "HTTP/1.1 400 Bad Request" + crlf
	case InternalServerError:
		reasonPhrase = "HTTP/1.1 500 Internal Server Error" + crlf
	default:
	}
	w.Write([]byte(reasonPhrase))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {

	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for header, val := range(headers) {

		headerPhrase := fmt.Sprintf("%s:%s%s", header, val, crlf)
		w.Write([]byte(headerPhrase))
	}

	return nil
}