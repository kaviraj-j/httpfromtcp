package main

import (
	"httpfromtcp/kaviraj-j/internal/request"
	"httpfromtcp/kaviraj-j/internal/response"
	"httpfromtcp/kaviraj-j/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func Handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/bad_request":
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "you fucked up something",
		}
	case "/server_error":
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "oops i'm sorry",
		}
	default:
		return nil
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
