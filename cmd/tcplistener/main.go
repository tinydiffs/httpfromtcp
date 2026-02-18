package main

import (
	// "errors"
	"fmt"
	// "io"
	"log"
	"net"
	// "strings"

	"github.com/tinydiffs/httpfromtcp/internal/request"
)

func main() {
	// file, err := os.Open("messages.txt")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("a connection has be accepted")

		req, err := request.RequestFromReader(conn)

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s",req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		// outputChan := getLinesChannel(conn)
		// for line := range(outputChan) {
		// 	fmt.Printf("%s\n", line)
		// }

		// fmt.Println("the connection has been closed")
	}

	// outputChan := getLinesChannel(file)

	// for line := range(outputChan) {
	// 	fmt.Printf("read: %s\n", line)
	// }

	// os.Exit(0)
	
}

// func getLinesChannel(f io.ReadCloser) <-chan string {

// 	var readingChan = make(chan string)

// 	go func() {
// 		defer close(readingChan)
// 		defer f.Close()
// 		b := make([]byte, 8)
// 		var currentLine []string

// 		for ; ; {
			
// 			readByte, err := f.Read(b)
// 			sections := strings.Split(string(b[:readByte]), "\n")
// 			if len(sections) > 1 {
// 				line := append(currentLine, strings.Join(sections[:len(sections) - 1], "\n"))
// 				readingChan <- strings.Join(line, "")

// 				currentLine = currentLine[:0]
// 				last := sections[len(sections) - 1]
// 				if last != "" {
// 					currentLine = append(currentLine, last)
// 				}
// 				continue
// 			}

// 			if sections[0] != "" {
// 				currentLine = append(currentLine, sections[0])
// 			}
			
// 			if err == nil {
// 				continue
// 			}

// 			if errors.Is(err, io.EOF) {
// 				if len(currentLine) > 0 {
// 					readingChan <- strings.Join(currentLine, "")
// 				}
// 				log.Println("EOF")
// 				return
// 			} else {
// 				log.Println(err)
// 				return
// 			}
// 		}
// 	}()

// 	return readingChan
// }