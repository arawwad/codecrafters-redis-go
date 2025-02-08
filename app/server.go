package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

var db = map[types.RespType]types.RespType{}

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

		instruction, args, ok := parseCommand(buf[:n])
		if !ok {
			break
		}

		println(instruction)

		switch instruction {
		case "PING":
			writeResponse(conn, types.SimpleString("PONG").Marshal())

		case "ECHO":
			writeResponse(conn, args[0].Marshal())

		case "SET":
			db[args[0]] = args[1]
			writeResponse(conn, types.SimpleString("OK").Marshal())
		case "GET":
			val, ok := db[args[0]]
			if !ok {
				writeResponse(conn, []byte(types.NullBulkString))
			} else {
				writeResponse(conn, val.Marshal())
			}
		default:
			fmt.Println("Error parsing command: unsupported command")
			break
		}

	}
}

func parseCommand(input []byte) (string, []types.RespType, bool) {
	p := parser.New(input)
	resp, ok := p.Parse()
	if !ok {
		fmt.Println("Error parsing command: ", p.Error.Error())
		return "", nil, false
	}

	a, ok := resp.(types.Array)
	if !ok {
		fmt.Println("Error parsing command: input is not valid Array")
		return "", nil, false
	}

	instruction, ok := a[0].(types.BulkString)
	if !ok {
		fmt.Println("Error parsing command: Command is not valid BulkString")
		return "", nil, false
	}

	return strings.ToUpper(string(instruction)), a[1:], true
}

func writeResponse(conn net.Conn, resp []byte) bool {
	_, writeErr := conn.Write(resp)
	if writeErr != nil {
		fmt.Println("Error writing response:", writeErr)
		return false
	}
	return true
}
