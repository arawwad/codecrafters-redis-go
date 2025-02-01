package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting server on port 6379...")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379:", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn) // Handle each connection in a goroutine
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error reading request:", err)
			}
			break
		}

		fmt.Println("Received:", string(buf[:n])) // Log received data
		_, writeErr := conn.Write([]byte("+PONG\r\n"))
		if writeErr != nil {
			fmt.Println("Error writing response:", writeErr)
			break
		}
	}
}
