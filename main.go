package main

import (
	"fluxKV/server"
	"fluxKV/utils"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	s := server.NewServer()

	if err = s.Listen(addr); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}
