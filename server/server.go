package server

import (
	"bufio"
	"fluxKV/commands"
	"fluxKV/internal"
	"fluxKV/utils"
	"fmt"
	"io"
	"net"
)

type Request struct {
	Conn net.Conn
	Cmd  []string
}
type Server struct {
	db       *internal.DataStore
	requests chan Request
}

func NewServer() *Server {
	return &Server{
		db:       internal.NewDataStore(),
		requests: make(chan Request, 1024),
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		cmd, err := internal.ReadRespCommand(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}
		s.requests <- Request{conn, cmd}
	}
}

func (s *Server) execute(req Request) {
	conn := req.Conn
	cmd := req.Cmd

	commands.ExecuteCommand(conn, cmd, s.db)
}

func (s *Server) eventLoop() {
	for req := range s.requests {
		s.execute(req)
	}
}

func (s *Server) Listen(addr *net.TCPAddr) error {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	utils.Logger.Info("Listening on %s", addr.String())

	go s.eventLoop()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}
