package client

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/db"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

const (
	ok     = types.SimpleString("OK")
	pong   = types.SimpleString("PONG")
	queued = types.SimpleString("QUEUED")
)

type Client struct {
	conn net.Conn
	*db.DB
	transactionMode bool
	queue           []*command
}

func New(db *db.DB, conn net.Conn) *Client {
	return &Client{
		DB:    db,
		conn:  conn,
		queue: []*command{},
	}
}

func (c *Client) respond(resp types.RespType) {
	_, writeErr := c.conn.Write(resp.Marshal())
	if writeErr != nil {
		fmt.Println("Error writing response:", writeErr)
	}
}

func (c *Client) Queue(cmd *command) {
	c.queue = append(c.queue, cmd)
	c.respond(queued)
}

func (c *Client) HandleConnection() {
	defer c.conn.Close()

	buf := make([]byte, 128)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error reading request:", err)
			}
			break
		}
		cmd, ok := parseCommand(buf[:n])
		if !ok {
			break
		}

		_, isExec := cmd.(execCmd)
		_, isDiscard := cmd.(discardCmd)

		shouldContinueTransaction := !(isExec || isDiscard)

		if c.transactionMode && shouldContinueTransaction {
			c.Queue(&cmd)
		} else {
			resp := cmd.exec(c)
			c.respond(resp)
		}

	}
}
