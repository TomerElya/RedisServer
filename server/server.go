package server

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

type Server struct {
	address  string
	port     int
	listener net.Listener
	store    store
	reqParser
}

func CreateServer(address string, port int) Server {
	return Server{address: address, port: port, store: CreateStore(), reqParser: createRequestParser()}
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
	logger := log.WithField("address", conn.RemoteAddr().String())
	logger.Info("new connection established")
	reader := bufio.NewReader(conn)
	var err error = nil
	var req request
	for err == nil {
		req, err = s.reqParser.constructRequest(reader)
		if err != nil {

		}
	}
	logger.WithError(err).Error("error received while listening to connection")
	_, err = conn.Write([]byte(err.Error()))
	if err != nil {
		logger.WithError(err).Error("failed to write error to client")
	}
	err = conn.Close()
	if err != nil {
		logger.WithError(err).Error("error while trying to close connection")
	}
}
