package commands

type SetCommand struct{}

func (s *SetCommand) execute(ctx *ctxInterface, cmd []string) {
	if len(cmd) != 3 {
		ctx.Conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	ctx.DB.Set(cmd[1], cmd[2])
	ctx.Conn.Write([]byte("+OK\r\n"))
}
