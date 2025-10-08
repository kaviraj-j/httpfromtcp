package main

import (
	"fmt"
	"httpfromtcp/kaviraj-j/internal/request"
	"httpfromtcp/kaviraj-j/internal/response"
	"httpfromtcp/kaviraj-j/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func Handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/bad_request" {
		body := "you fucked up"
		w.WriteStatusLine(400)
		h := response.GetDefaultHeaders(len(body))
		h.OverrideValue("Key", "MyVal")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	} else if req.RequestLine.RequestTarget == "/server_error" {
		body := "sorry my bad"
		w.WriteStatusLine(500)
		h := response.GetDefaultHeaders(len(body))
		h.OverrideValue("Key", "MyVal")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
		target := req.RequestLine.RequestTarget
		url := "https://httpbin.org/" + target[len("/httpbin/"):]
		res, err := http.Get(url)
		if err != nil {
			w.WriteStatusLine(500)
			return
		}
		w.WriteStatusLine(response.StatusOk)
		h := response.GetDefaultHeaders(0)
		h.Delete("Content-Length")
		h.OverrideValue("transfer-encoding", "chunked")
		h.OverrideValue("content-type", "text/plain")
		w.WriteHeaders(h)
		for {
			data := make([]byte, 128)
			n, err := res.Body.Read(data)
			if err != nil {
				break
			}

			w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
			w.WriteBody(data[:n])
			w.WriteBody([]byte("\r\n"))
		}
	} else if req.RequestLine.RequestTarget == "/video" {
		fileData, _ := os.ReadFile("./assets/earth.mp4")
		h := response.GetDefaultHeaders(len(fileData))
		h.OverrideValue("Content-Type", "video/mp4")
		w.WriteStatusLine(response.StatusOk)
		w.WriteHeaders(h)
		w.WriteBody(fileData)
	} else {
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
