package client

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/db"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type Client struct {
	conn net.Conn
	*db.DB
	transactionMode bool
	queue           []*Command
}

func New(db *db.DB, conn net.Conn) *Client {
	return &Client{
		DB:    db,
		conn:  conn,
		queue: []*Command{},
	}
}

func (c *Client) Respond(resp types.RespType) {
	_, writeErr := c.conn.Write(resp.Marshal())
	if writeErr != nil {
		fmt.Println("Error writing response:", writeErr)
	}
}

func (c *Client) OK() {
	c.Respond(types.SimpleString("OK"))
}

func (c *Client) PONG() {
	c.Respond(types.SimpleString("PONG"))
}

func (c *Client) Queue(cmd *Command) {
	c.queue = append(c.queue, cmd)
	c.Respond(types.SimpleString("QUEUED"))
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
		cmd, ok := ParseCommand(buf[:n])
		if !ok {
			break
		}

		_, isExec := cmd.(Exec)

		if c.transactionMode && !isExec {
			c.Queue(&cmd)
		} else {
			resp := cmd.Exec(c)
			c.Respond(resp)
		}

	}
}
