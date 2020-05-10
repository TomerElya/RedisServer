package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

type Server struct {
	address  string
	port     int
	listener net.Listener
	store    store
}

func CreateServer(address string, port int) Server {
	return Server{address: address, port: port, store: CreateStore()}
}

func (s *Server) StartAndListen() {
	address := fmt.Sprintf("%s:%d", s.address, s.port)
	listener, err := net.Listen("tcp", address)
	s.listener = listener
	if err != nil {
		log.WithError(err).WithField("address", address).Error("listener failed to start")
		panic(err)
	}
	log.WithField("address", address).Info("listener successfully started")
	s.listen()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.WithError(err).Error("failed receiving new connection")
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log.WithField("address", conn.RemoteAddr()).Info("new connection established")
}
