package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/tinydiffs/httpfromtcp/internal/request"
	"github.com/tinydiffs/httpfromtcp/internal/response"
)

type Server struct{

	listener	net.Listener
	active		atomic.Bool
	handler		Handler

}

// type HandlerError struct{

// 	StatusCode		response.StatusCode
// 	Message			string
// }

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {

	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return &Server{}, err
	}

	server := Server{
		listener: listener,
		handler: handler,
	}

	go server.listen()
	return &server, nil
}



func (s *Server) Close() error {
	s.active.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}



func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.active.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}


func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	writer := response.Writer{
		Connection: conn,
	}

	req, err := request.RequestFromReader(conn)
	if err != nil {
		innerErr := writer.WriteStatusLine(response.BadRequest)
		if innerErr != nil {
			log.Printf("error writing status line: %s", innerErr)
		}

		body := []byte(fmt.Sprintf("<html><body><h1>400 Bad Request</h1><p>%v</p></body></html>", err))
		innerErr = writer.WriteHeaders(response.GetDefaultHeaders(len(body)))
		if innerErr != nil {
			log.Printf("error writing headers: %s", innerErr)
		}

		_, innerErr = writer.WriteBody(body)
		if innerErr != nil {
			log.Printf("error writing body: %s", innerErr)
		}

		return
	}

	// buf := bytes.Buffer{}
	

	s.handler(&writer, req)
	// if hErr != nil {
	// 	err = hErr.Write(conn)
	// 	if err != nil {
	// 		log.Printf("error handling(writing) error: %s", err)
	// 	}
	// 	return
	// }

	// defaultHeaders := response.GetDefaultHeaders(buf.Len())
	// err = response.WriteStatusLine(conn, response.Ok)
	// if err != nil {
	// 	log.Printf("error writing status line: %s", err)

	// }

	// err = response.WriteHeaders(conn, defaultHeaders)
	// if err != nil {
	// 	log.Printf("error writing headers: %s", err)
	// }

	// _, err = conn.Write(buf.Bytes())
	// if err != nil {
	// 	log.Printf("error writing budy: %s", err)
	// }

	fmt.Println("Responded!")
}

// func(h *HandlerError) Write(w io.Writer) error {

// 	err := response.WriteStatusLine(w, h.StatusCode)
// 	if err != nil {
// 		return err
// 	}

// 	defaultHeaders := response.GetDefaultHeaders(len(h.Message))
// 	err = response.WriteHeaders(w, defaultHeaders)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = w.Write([]byte(h.Message))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }