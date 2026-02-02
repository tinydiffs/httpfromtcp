package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	outputChan := getLinesChannel(file)

	for line := range(outputChan) {
		fmt.Printf("read: %s\n", line)
	}

	os.Exit(0)
	

	// b := make([]byte, 8)
	// var currentLine []string
	// var i int

	// for ; ; {
		
	// 	readByte, err := file.ReadAt(b, int64(i))
	// 	i = i + readByte
	// 	sections := strings.Split(string(b[:readByte]), "\n")
	// 	if len(sections) > 1 {
	// 		line := append(currentLine, strings.Join(sections[:len(sections) - 1], "\n"))
	// 		fmt.Printf("read: %s\n", strings.Join(line, ""))

	// 		currentLine = currentLine[:0]
	// 		last := sections[len(sections) - 1]
	// 		if last != "" {
	// 			currentLine = append(currentLine, last)
	// 		}
	// 		continue
	// 	}

	// 	if sections[0] != "" {
	// 		currentLine = append(currentLine, sections[0])
	// 	}
		
	// 	if err == nil {
	// 		continue
	// 	}

	// 	if errors.Is(err, io.EOF) {
	// 		if len(currentLine) > 0 {
	// 			fmt.Printf("read: %s\n", strings.Join(currentLine, ""))
	// 			os.Exit(0)
	// 		}
	// 		os.Exit(0)
	// 	} else {
	// 		log.Fatal("err")
	// 	}
	// }
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	var readingChan = make(chan string)

	go func() {
		defer close(readingChan)
		defer f.Close()
		b := make([]byte, 8)
		var currentLine []string

		for ; ; {
			
			readByte, err := f.Read(b)
			sections := strings.Split(string(b[:readByte]), "\n")
			if len(sections) > 1 {
				line := append(currentLine, strings.Join(sections[:len(sections) - 1], "\n"))
				readingChan <- strings.Join(line, "")
				// fmt.Printf("read: %s\n", strings.Join(line, ""))

				currentLine = currentLine[:0]
				last := sections[len(sections) - 1]
				if last != "" {
					currentLine = append(currentLine, last)
				}
				continue
			}

			if sections[0] != "" {
				currentLine = append(currentLine, sections[0])
			}
			
			if err == nil {
				continue
			}

			if errors.Is(err, io.EOF) {
				if len(currentLine) > 0 {
					readingChan <- strings.Join(currentLine, "")
					// fmt.Printf("read: %s\n", strings.Join(currentLine, ""))
				}
				log.Println("EOF")
				return
			} else {
				log.Println(err)
				return
			}
		}
	}()

	return readingChan
}