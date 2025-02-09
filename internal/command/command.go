package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/parser"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type Command interface {
	name() string
}

func Parse(input []byte) (Command, bool) {
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
	}

	return nil, false
}

type Ping struct{}

func (Ping) name() string {
	return "ping"
}

type Echo struct {
	Value types.RespType
}

func (Echo) name() string {
	return "echo"
}

type Get struct {
	Key types.RespType
}

func (Get) name() string {
	return "get"
}

type Set struct {
	Key   types.RespType
	Value types.RespType
	TTL   *time.Duration
}

func (Set) name() string {
	return "set"
}

type Incr struct {
	Key types.RespType
}

func (Incr) name() string {
	return "incr"
}

type Multi struct{}

func (Multi) name() string {
	return "incr"
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
