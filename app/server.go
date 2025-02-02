package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/marshal"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/parser"
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

		p := parser.New(buf[:n])
		cmd, ok := p.ParseArray()
		if !ok {
			fmt.Println("Error parsing command: ", p.Error.Error())
			break
		}

		switch strings.ToUpper(cmd[0]) {
		case "PING":
			writeResponse(conn, []byte("+PONG\r\n"))

		case "ECHO":
			writeResponse(conn, marshal.MarshalBulkString(cmd[1]))
		}
	}
}

func writeResponse(conn net.Conn, resp []byte) bool {
	_, writeErr := conn.Write(resp)
	if writeErr != nil {
		fmt.Println("Error writing response:", writeErr)
		return false
	}
	return true
}
