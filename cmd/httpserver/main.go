package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tinydiffs/httpfromtcp/internal/headers"
	"github.com/tinydiffs/httpfromtcp/internal/request"
	"github.com/tinydiffs/httpfromtcp/internal/response"
	"github.com/tinydiffs/httpfromtcp/internal/server"
)



const port = 42069

func main() {
	serv, err := server.Serve(port, httpBinProxy)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer serv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func aHandler(w *response.Writer,r *request.Request) {

	switch r.RequestLine.RequestTarget{
	case "/yourproblem":
		w.WriteStatusLine(response.BadRequest)
		body := []byte(
			`<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
			</html>`)

		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
	case "/myproblem":
		w.WriteStatusLine(response.InternalServerError)
		body := []byte(
			`<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
			</html>`)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
	default:
		w.WriteStatusLine(response.Ok)
		body := []byte(
			`<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
			</html>`)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
	}
}

func httpBinProxy(w *response.Writer, r *request.Request) {

	urlPrefix := "https://httpbin.org/"

	if !strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
		log.Printf("Incorrect request for handler")
		return
	}

	requestPath := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	resp, err := http.Get(urlPrefix + requestPath)
	if err != nil {
		log.Print(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("httpbin response error: %s", resp.Status)
		return
	}
	w.WriteStatusLine(response.Ok)
	headers := headers.NewHeaders()
	headers["content-type"] = "text/plain"
	headers["transfer-encoding"] = "chunked"
	w.WriteHeaders(headers)

	buf := make([]byte, 1024)
	for {
		readBytes, err := resp.Body.Read(buf)
		if readBytes > 0 {
			wroteBytes, _ := w.WriteChunkedBody(buf[:readBytes])
			log.Printf("Read and Wrote: %d bytes", wroteBytes)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			break
		}

	}

	w.WriteChunkedBodyDone()
}