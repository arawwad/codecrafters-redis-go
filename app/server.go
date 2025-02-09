package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/db"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type dbValue struct {
	value   types.RespType
	expires time.Time
}

var DB = db.New()

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

		cmd, ok := command.Parse(buf[:n])

		if !ok {
			break
		}

		switch cmd.(type) {
		case command.Ping:
			writeResponse(conn, types.SimpleString("PONG").Marshal())

		case command.Echo:
			writeResponse(conn, cmd.(command.Echo).Value.Marshal())

		case command.Set:
			setCmd := cmd.(command.Set)
			DB.Set(setCmd.Key, setCmd.Value, setCmd.TTL)
			writeResponse(conn, types.SimpleString("OK").Marshal())

		case command.Get:
			val, ok := DB.Get(cmd.(command.Get).Key)
			if !ok {
				writeResponse(conn, []byte(types.NullBulkString))
			}
			writeResponse(conn, val.Marshal())

		case command.Incr:
			writeResponse(conn, DB.Incr(cmd.(command.Incr).Key).Marshal())

		default:
			fmt.Println("Error parsing command: unsupported command")
			break
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
