package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
)

type Server struct {
	address             string
	port                int
	listener            net.Listener
	reqParser           RequestsParser
	cmdHandler          CommandHandler
	acceptedConnections chan net.Conn
	signalChannel       chan os.Signal
}

func CreateServer(address string, port int) Server {
	return Server{
		address:             address,
		port:                port,
		cmdHandler:          CreateCommandHandler(),
		reqParser:           createRequestParser(),
		acceptedConnections: make(chan net.Conn),
		signalChannel:       make(chan os.Signal),
	}
}

func (s *Server) StartAndListen() {
	address := fmt.Sprintf("%s:%d", s.address, s.port)
	listener, err := net.Listen("tcp", address)
	signal.Notify(s.signalChannel, os.Interrupt)
	s.listener = listener
	if err != nil {
		log.WithError(err).WithField("address", address).Error("listener failed to start")
		panic(err)
	}
	log.WithField("address", address).Info("listener successfully started")
	go s.listenForConnections()
	s.run()
}

func (s *Server) run() {
	for {
		select {
		case conn := <-s.acceptedConnections:
			go s.handleConnection(conn)
		case <-s.signalChannel:
			log.Info("interrupt received from console, exiting...")
			return
		}
	}
}

func (s *Server) listenForConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.WithError(err).Error("failed receiving new connection")
		} else {
			s.acceptedConnections <- conn
		}

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	client := CreateClient(conn)
	var err error = nil
	var req Request
	for err == nil {
		req, err = s.reqParser.ConstructRequest(client.reader)
		if err == nil {
			req.client = client
			go s.cmdHandler.AppendRequest(req)
		}
	}
	client.Disconnect(err)
}
