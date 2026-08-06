package commands

import (
	"fluxKV/internal"
	"net"
)

type SetCommand struct{}

func (s *SetCommand) execute(conn net.Conn, cmd []string, db *internal.DataStore) {
	if len(cmd) != 3 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	db.Set(cmd[1], cmd[2])
	conn.Write([]byte("+OK\r\n"))
}
