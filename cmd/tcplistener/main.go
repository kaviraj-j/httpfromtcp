package main

import (
	"bytes"
	"errors"
	"fmt"
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
			linesChan := getLinesChannel(newConn)
			for l := range linesChan {
				fmt.Println(l)
			}
			fmt.Println("connection is closed")
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
