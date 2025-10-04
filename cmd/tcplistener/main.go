package main

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/kaviraj-j/internal/request"
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("waiting to accept connection")
	for {
		newConn, err := conn.Accept()
		if err != nil {
			fmt.Println("error while accepting a connection")
			newConn.Close()
			continue
		}

		go func() {
			defer newConn.Close()
			fmt.Println("new connection is accepted")
			r, _ := request.RequestFromReader(newConn)
			fmt.Println("====== Request created =======")
			fmt.Printf("Method: %s\nHttp Version: %s\nEndpoint: %s\n", r.RequestLine.Method, r.RequestLine.HttpVersion, r.RequestLine.RequestTarget)
			// print headers
			fmt.Println("--- Headers---")
			for k, v := range r.RequestHeaders {
				fmt.Printf("%s: %s\n", k, v)
			}
			fmt.Println("---- Body ----")
			fmt.Println(string(r.Body))
		}()
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)
	go func() {
		line := ""
		defer close(out)
		for {
			data := make([]byte, 8)

			n, err := f.Read(data)
			if err != nil {
				if errors.Is(err, io.EOF) && len(line) > 0 {
					out <- line
				}
				break
			}
			data = data[:n]
			if idx := bytes.IndexByte(data, '\n'); idx != -1 {
				line += string(data[:idx])
				data = data[idx+1:]
				out <- line
				line = ""
			}
			line += string(data)

		}
	}()
	return out
}
