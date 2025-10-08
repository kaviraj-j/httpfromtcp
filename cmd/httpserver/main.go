package main

import (
	"fmt"
	"httpfromtcp/kaviraj-j/internal/request"
	"httpfromtcp/kaviraj-j/internal/response"
	"httpfromtcp/kaviraj-j/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func Handler(w *response.Writer, req *request.Request) {
	fmt.Println(req.RequestLine.RequestTarget)
	switch req.RequestLine.RequestTarget {
	case "/bad_request":
		body := "you fucked up"
		w.WriteStatusLine(400)
		h := response.GetDefaultHeaders(len(body))
		h.OverrideValue("Key", "MyVal")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	case "/server_error":
		body := "sorry my bad"
		w.WriteStatusLine(500)
		h := response.GetDefaultHeaders(len(body))
		h.OverrideValue("Key", "MyVal")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	default:
		body := "Status OK"
		w.WriteStatusLine(200)
		h := response.GetDefaultHeaders(len(body))
		h.OverrideValue("Key", "MyVal")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	}
}
func main() {
	server, err := server.Serve(port, Handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
