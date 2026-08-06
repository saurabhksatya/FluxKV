package commands

import (
	"fluxKV/internal"
	"net"
	"strings"
)

type CommandInterface interface {
	execute(conn net.Conn, cmd []string, db *internal.DataStore)
}

var Commands = map[string]CommandInterface{}

func register(name string, command CommandInterface) {
	Commands[name] = command
}

func init() {
	register("PING", &PingCommand{})
	register("GET", &GetCommand{})
	register("SET", &SetCommand{})
}

func ExecuteCommand(conn net.Conn, cmd []string, db *internal.DataStore) {
	if len(cmd) == 0 {
		conn.Write([]byte("-ERR empty command\r\n"))
		return
	}
	finalCmd := strings.ToUpper(cmd[0])

	inter, ok := Commands[finalCmd]
	if !ok {
		conn.Write([]byte("-ERR unknown command\r\n"))
		return
	}

	inter.execute(conn, cmd, db)
}
