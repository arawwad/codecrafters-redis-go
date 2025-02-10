package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type command interface {
	exec(*Client) types.RespType
}

func parseCommand(input []byte) (command, bool) {
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
		return pingCmd{}, true
	case "ECHO":
		return echoCmd{Value: arr[1]}, true
	case "GET":
		return getCmd{Key: arr[1]}, true
	case "SET":
		return setCmd{Key: arr[1], Value: arr[2], TTL: getTTL(arr)}, true
	case "INCR":
		return incrCmd{Key: arr[1]}, true
	case "MULTI":
		return multiCmd{}, true
	case "EXEC":
		return execCmd{}, true
	case "DISCARD":
		return discardCmd{}, true
	case "TYPE":
		return typeCmd{Key: arr[1]}, true
	case "XADD":
		return xAddCmd{key: arr[1], id: arr[2], entryKey: arr[3], entryValue: arr[4]}, true
	}

	return nil, false
}

type pingCmd struct{}

func (pingCmd) exec(c *Client) types.RespType {
	return pong
}

type echoCmd struct {
	Value types.RespType
}

func (e echoCmd) exec(c *Client) types.RespType {
	return e.Value
}

type getCmd struct {
	Key types.RespType
}

func (cmd getCmd) exec(c *Client) types.RespType {
	val, ok := c.Get(cmd.Key)
	if !ok {
		return types.NullBulkString{}
	}
	return val
}

type setCmd struct {
	Key   types.RespType
	Value types.RespType
	TTL   *time.Duration
}

func (cmd setCmd) exec(c *Client) types.RespType {
	c.Set(cmd.Key, cmd.Value, cmd.TTL)
	return ok
}

type incrCmd struct {
	Key types.RespType
}

func (cmd incrCmd) exec(c *Client) types.RespType {
	return c.Incr(cmd.Key)
}

type typeCmd struct {
	Key types.RespType
}

func (cmd typeCmd) exec(c *Client) types.RespType {
	val, ok := c.Get(cmd.Key)
	if !ok {
		return types.NoneType
	}
	return val.Type()
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
