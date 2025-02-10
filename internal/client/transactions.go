package client

import "github.com/codecrafters-io/redis-starter-go/internal/resp/types"

type multiCmd struct{}

func (multiCmd) exec(c *Client) types.RespType {
	c.transactionMode = true
	return ok
}

type execCmd struct{}

func (execCmd) exec(c *Client) types.RespType {
	if !c.transactionMode {
		return types.SimpleError("ERR EXEC without MULTI")
	}
	c.transactionMode = false

	result := []types.RespType{}
	for _, cmd := range c.queue {
		result = append(result, (*cmd).exec(c))
	}

	return types.Array(result)
}

type discardCmd struct{}

func (discardCmd) exec(c *Client) types.RespType {
	if !c.transactionMode {
		return types.SimpleError("ERR DISCARD without MULTI")
	}
	c.transactionMode = false
	c.queue = []*command{}

	return ok
}
