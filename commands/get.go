package commands

import (
	"fmt"
)

type GetCommand struct{}

func (s *GetCommand) execute(ctx *ctxInterface, cmd []string) {
	if len(cmd) != 2 {
		ctx.Conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	value, ok := ctx.DB.Get(cmd[1])
	if !ok {
		ctx.Conn.Write([]byte("$-1\r\n"))
		return
	}

	fmt.Fprintf(ctx.Conn, "$%d\r\n%s\r\n", len(value), value)
}
