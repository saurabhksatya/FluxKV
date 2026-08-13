package server

import (
	"bufio"
	"fluxKV/commands"
	"fluxKV/configuration"
	"fluxKV/internal"
	"fluxKV/replication"
	"fluxKV/utils"
	"fmt"
	"io"
	"net"
	"strings"
)

type Request struct {
	Conn net.Conn
	Cmd  []string
}
type Server struct {
	db       *internal.DataStore
	requests chan Request

	self *configuration.Server
}

func NewServer(cfg *configuration.Config, self *configuration.Server) *Server {
	s := &Server{
		requests: make(chan Request, 1024),
		self:     self,
	}

	switch self.Role {
	case "main":
		s.db = internal.NewDataStore(true)
		s.db.SetReplicationServer(replication.NewReplicationServer(s.db))
	case "replica":
		s.db = internal.NewDataStore(false)
		s.db.SetReplicationClient(replication.NewReplicationClient(self.MainAddr, self.ID, s.db))
	}

	return s
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

	if len(cmd) == 0 {
		return
	}

	name := strings.ToUpper(cmd[0])

	if commands.IsWriteCommand(name) && s.self.Role == "replica" {
		conn.Write([]byte("-ERR READONLY You can't write against a replica.\r\n"))
		return
	}

	commands.ExecuteCommand(conn, cmd, s.db)
}

func (s *Server) eventLoop() {
	for req := range s.requests {
		s.execute(req)
	}
}

func (s *Server) Listen() error {
	s.db.StartReplication(s.self.GRPCListen)

	addr, err := net.ResolveTCPAddr("tcp", s.self.Listen)
	if err != nil {
		return err
	}

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
