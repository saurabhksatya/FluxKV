package commands

import (
	"fluxKV/internal"
	"net"
)

type PingCommand struct{}

func (p *PingCommand) execute(conn net.Conn, cmd []string, db *internal.DataStore) {
	conn.Write([]byte("+PONG\r\n"))
}
