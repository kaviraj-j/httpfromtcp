package server

import (
	"fmt"
	"httpfromtcp/kaviraj-j/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
	}
	go server.listen()
	return server, err
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("accept error:", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	if s.closed.Load() {
		return nil
	}
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// body := "Hello, world!"

	// _, err := conn.Write([]byte(response))

	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		fmt.Println("write error:", err)
	}
	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		fmt.Println("write error:", err)
	}
}
