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
	headers["content-type"] = "text/html"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for header, val := range(headers) {

		headerPhrase := fmt.Sprintf("%s:%s%s", header, val, crlf)
		w.Write([]byte(headerPhrase))
	}
	w.Write([]byte(crlf))

	return nil
}

type writerState int

const(
	statusLine writerState = iota
	header
	body
)

type Writer struct{
	Connection	io.Writer
	writerState	writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writerState != statusLine {
		return fmt.Errorf("statusline must be written first")
	}

	err := WriteStatusLine(w.Connection, statusCode)
	if err != nil {
		return err
	}
	w.writerState = header
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.writerState != header {
		return fmt.Errorf("headers must be written after statusline")
	}

	err := WriteHeaders(w.Connection, headers)
	if err != nil {
		return err
	}
	w.writerState = body
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.writerState != body {
		return 0, fmt.Errorf("body must be wrtten after headers")
	}

	return w.Connection.Write(p)
}