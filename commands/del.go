package commands

import (
	"fmt"
)

type DelCommand struct{}

func (d *DelCommand) execute(ctx *ctxInterface, cmd []string) {
	if len(cmd) < 2 {
		ctx.Conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
		return
	}

	deleted := 0

	for _, key := range cmd[1:] {
		if ctx.DB.Delete(key) {
			deleted++
		}
	}

	fmt.Fprintf(ctx.Conn, ":%d\r\n", deleted)
}
