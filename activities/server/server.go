// server.go
package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

type Server struct {
	clients map[net.Conn]*Client
	mu      sync.Mutex
	counter int
}

func (s *Server) handleConnection(conn net.Conn) {
	s.mu.Lock()
	s.counter++
	name := fmt.Sprintf("Anonymous %d", s.counter)
	s.clients[conn] = &Client{conn: conn, name: name}
	s.mu.Unlock()

	fmt.Printf("%s joined the chat\n", name)
	s.broadcast(conn, fmt.Sprintf("%s joined the chat\n", name))

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		fmt.Printf("%s left the chat\n", name)
		s.broadcast(conn, fmt.Sprintf("%s left the chat\n", name))
		s.checkEmpty()
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return // Ngắt kết nối khi có lỗi đọc
		}
		s.broadcast(conn, fmt.Sprintf("%s: %s", name, string(buf[:n])))
	}
}

func (s *Server) broadcast(sender net.Conn, msg string) {
	s.mu.Lock()
	for _, client := range s.clients {
		if client.conn != sender {
			_, err := client.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	s.mu.Unlock()
}

func (s *Server) checkEmpty() {
	s.mu.Lock()
	if len(s.clients) == 0 {
		fmt.Println("No clients left. Shutting down server...")
		os.Exit(0)
	}
	s.mu.Unlock()
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	fmt.Println("Server listening on port 8080")

	server := &Server{clients: make(map[net.Conn]*Client)}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go server.handleConnection(conn)
	}
}
