package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type Command interface {
	Exec(*Client) types.RespType
}

func ParseCommand(input []byte) (Command, bool) {
	p := parser.New(input)
	resp, ok := p.Parse()
	if !ok {
		fmt.Println("Error parsing command: ", p.Error.Error())
		return nil, false
	}

	arr, ok := resp.(types.Array)
	if !ok {
		fmt.Println("Error parsing command: input is not valid Array")
		return nil, false
	}

	instruction, ok := arr[0].(types.BulkString)
	if !ok {
		fmt.Println("Error parsing command: Command is not valid BulkString")
		return nil, false
	}

	switch strings.ToUpper(string(instruction)) {
	case "PING":
		return Ping{}, true
	case "ECHO":
		return Echo{Value: arr[1]}, true
	case "GET":
		return Get{Key: arr[1]}, true
	case "SET":
		return Set{Key: arr[1], Value: arr[2], TTL: getTTL(arr)}, true
	case "INCR":
		return Incr{Key: arr[1]}, true
	case "MULTI":
		return Multi{}, true
	case "EXEC":
		return Exec{}, true
	case "DISCARD":
		return Discard{}, true
	}

	return nil, false
}

type Ping struct{}

func (Ping) Exec(c *Client) types.RespType {
	return types.SimpleString("PONG")
}

type Echo struct {
	Value types.RespType
}

func (e Echo) Exec(c *Client) types.RespType {
	return e.Value
}

type Get struct {
	Key types.RespType
}

func (cmd Get) Exec(c *Client) types.RespType {
	val, ok := c.Get(cmd.Key)
	if !ok {
		return types.NullBulkString{}
	}
	return val
}

type Set struct {
	Key   types.RespType
	Value types.RespType
	TTL   *time.Duration
}

func (cmd Set) Exec(c *Client) types.RespType {
	c.Set(cmd.Key, cmd.Value, cmd.TTL)
	return types.SimpleString("OK")
}

type Incr struct {
	Key types.RespType
}

func (cmd Incr) Exec(c *Client) types.RespType {
	return c.Incr(cmd.Key)
}

type Multi struct{}

func (Multi) Exec(c *Client) types.RespType {
	c.transactionMode = true
	return types.SimpleString("OK")
}

type Exec struct{}

func (Exec) Exec(c *Client) types.RespType {
	if !c.transactionMode {
		return types.SimpleError("ERR EXEC without MULTI")
	}
	c.transactionMode = false

	result := []types.RespType{}
	for _, cmd := range c.queue {
		result = append(result, (*cmd).Exec(c))
	}

	return types.Array(result)
}

type Discard struct{}

func (Discard) Exec(c *Client) types.RespType {
	if !c.transactionMode {
		return types.SimpleError("ERR DISCARD without MULTI")
	}
	c.transactionMode = false
	c.queue = []*Command{}

	return types.SimpleString("OK")
}

func getTTL(args []types.RespType) *time.Duration {
	if len(args) < 3 {
		return nil
	}

	for index, value := range args {
		if str, ok := value.(types.BulkString); ok && strings.ToLower(string(str)) == "px" {
			if num, ok := args[index+1].Num(); ok {
				duration := time.Duration(num) * time.Millisecond
				return &duration
			} else {
				return nil
			}
		}
	}
	return nil
}
