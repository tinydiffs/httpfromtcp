package request

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/tinydiffs/httpfromtcp/internal/headers"
)

const bufferSize int = 8
const crlf string = "\r\n"

type state int

const(
	initialized state = iota
	parsingHeaders
	parsingBody
	done
)

type Request struct {
	readState		state
	RequestLine 	RequestLine
	Headers			headers.Headers
	Body			[]byte
	bodyLengthRead	int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buff := make([]byte, bufferSize)
	readToIndex := 0
	req := Request{
		readState: initialized,
		Headers: headers.NewHeaders(),
	}

	for req.readState != done {

		if readToIndex == len(buff) {
			newBuff := make([]byte, len(buff) * 2)
			copy(newBuff, buff)
			buff = newBuff
		}

		readbytes, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.readState, readbytes)
		}
		if err != nil {
			return &Request{}, fmt.Errorf("error reading file %v", err)
		}
		readToIndex += readbytes

		parsedBytes, err := req.parse(buff[:readToIndex])
		if err != nil {
			return &Request{}, fmt.Errorf("error parsing data %v", err)
		}

		copy(buff, buff[parsedBytes:])
		readToIndex -= parsedBytes
	}

	return &req, nil
}

func parseRequestLine(line []byte) (int, RequestLine, error) {

	lineString := string(line)
	parts := strings.Split(lineString, crlf)
	if len(parts) < 1 {
		return 0, RequestLine{}, fmt.Errorf("request parsing error")
	}

	// If no crlf was found
	if len(parts) < 2 {
		return 0, RequestLine{}, nil
	}

	reqLineParts := strings.Split(parts[0], " ")
	if len(reqLineParts) < 3 {
		return 0, RequestLine{}, fmt.Errorf("request line unknown format")
	}

	method := reqLineParts[0]
	requestTarget := reqLineParts[1]
	httpNameVersion := reqLineParts[2]

	for _, letter := range method {
		if !unicode.IsUpper(letter) {
			return 0, RequestLine{}, fmt.Errorf("wrong method format")
		}
	}

	httpNameVersionParts := strings.Split(httpNameVersion, "/")
	if len(httpNameVersionParts) < 2 {
		return 0,RequestLine{}, fmt.Errorf("unknown http version")
	}
	httpVersion := httpNameVersionParts[1]
	if httpVersion != "1.1" {
		return 0,RequestLine{}, fmt.Errorf("unsupported version: %s", httpVersion)
	}

	return len(parts[0]) + 2, RequestLine{
		HttpVersion: httpVersion,
		RequestTarget: requestTarget,
		Method: method,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {

	totalBytesParsed := 0
	for r.readState != done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil	
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.readState {
	case initialized:
		consumedBytes, reqline, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if consumedBytes == 0 {
			return 0, nil
		}

		r.RequestLine = reqline
		r.readState = parsingHeaders

		return consumedBytes, nil

	case parsingHeaders:
		consumedBytes, finished, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if finished {
			r.readState = parsingBody
			return consumedBytes, nil
		}
		return consumedBytes, nil
	
	case parsingBody:
		content_length_string := r.Headers.Get("Content-Length")

		if content_length_string == "" {
			r.readState = done
			return 0, nil
		}

		content_length, err := strconv.Atoi(content_length_string)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %s", err)
		}

		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)

		if r.bodyLengthRead > content_length {
			return 0, fmt.Errorf("Content-Length too large")
		}

		if r.bodyLengthRead == content_length {
			r.readState = done
			return len(data), nil
		}
		return len(data), nil

	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}