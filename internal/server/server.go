package server

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/tinydiffs/httpfromtcp/internal/response"
)

type Server struct{

	listener	net.Listener
	active		bool

}

func Serve(port int) (*Server, error) {

	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return &Server{}, err
	}

	server := Server{
		listener: listener,
		active: true,
	}

	go server.listen()
	return &server, nil
}



func (s *Server) Close() error {
	s.active = false
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}



func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go s.handle(conn)
	}
}


func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// body := "Hello World!"

	err := response.WriteStatusLine(conn, response.Ok)
	if err != nil {
		log.Printf("error writing status line: %s", err)
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("error writing headers: %s", err)
	}

	conn.Write([]byte("\r\n"))

	// response := fmt.Sprintf(
	// 	"HTTP/1.1 200 OK\r\n"+
	// 		"Content-Type: text/plain\r\n"+
	// 		// "Content-Length: %d\r\n"+
	// 		"\r\n"+
	// 		"%s",
	// 	// len(body),
	// 	body,
	// )

	// conn.Write([]byte(response))

	fmt.Println("Responded!")
}