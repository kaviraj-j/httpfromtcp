package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	line := ""
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
			fmt.Printf("read: %s\n", line)
			line = ""
		}
		line += string(data)

	}
}
