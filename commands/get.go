package commands

import (
	"fluxKV/internal"
	"fmt"
	"net"
)

type GetCommand struct{}

func (s *GetCommand) execute(conn net.Conn, cmd []string, db *internal.DataStore) {
	if len(cmd) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	value, ok := db.Get(cmd[1])
	if !ok {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
}
