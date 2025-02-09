package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/db"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type dbValue struct {
	value   types.RespType
	expires time.Time
}

func main() {
	fmt.Println("Starting server on port 6379...")

	DB := db.New()

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
		c := client.New(&DB, conn)
		go c.HandleConnection()
	}
}
