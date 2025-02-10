package client

import "github.com/codecrafters-io/redis-starter-go/internal/resp/types"

type xAddCmd struct {
	key        types.RespType
	id         types.RespType
	entryKey   types.RespType
	entryValue types.RespType
}

func (cmd xAddCmd) exec(c *Client) types.RespType {
	return c.AppendToStream(cmd.key, cmd.id, cmd.entryKey, cmd.entryValue)
}
