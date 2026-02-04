package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	udpCon, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpCon.Close()

	rder := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := rder.ReadString('\n')
		if err != nil {
			log.Println(err)
		}

		_, err = udpCon.Write([]byte(line))
		if err != nil {
			log.Println(err)
		}
	}
}