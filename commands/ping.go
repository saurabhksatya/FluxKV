package commands

type PingCommand struct{}

func (p *PingCommand) execute(ctx *ctxInterface, cmd []string) {
	ctx.Conn.Write([]byte("+PONG\r\n"))
}
