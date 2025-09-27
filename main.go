package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string)
	go func() {
		line := ""
		defer file.Close()
		defer close(out)
		for {

			data := make([]byte, 8)
			n, err := file.Read(data)
			if err != nil {
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
