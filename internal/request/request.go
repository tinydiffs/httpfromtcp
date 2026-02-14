package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

const bufferSize int = 8

type state int

const(
	initialized state = iota
	done
)

type Request struct {
	readState	  state
	RequestLine RequestLine
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
	}

	for req.readState != done {

		if readToIndex == len(buff) {
			newBuff := make([]byte, len(buff) * 2)
			copy(newBuff, buff)
			buff = newBuff
		}

		readbytes, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			req.readState = done
			break
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
	parts := strings.Split(lineString, "\r\n")
	if len(parts) < 1 {
		return 0, RequestLine{}, fmt.Errorf("request parsing error")
	}

	// If no \r\n was found
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

	switch r.readState {
	case initialized:
		consumedbytes, reqline, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if consumedbytes == 0 {
			return 0, nil
		}

		r.RequestLine = reqline
		r.readState = done

		return consumedbytes, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
	
}