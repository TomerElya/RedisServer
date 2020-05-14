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
	cmdHandler          CommandHandler
	acceptedConnections chan net.Conn
	signalChannel       chan os.Signal
	clientSet           map[string]Client
}

func CreateServer(address string, port int) Server {
	return Server{
		address:             address,
		port:                port,
		cmdHandler:          CreateCommandHandler(),
		acceptedConnections: make(chan net.Conn),
		signalChannel:       make(chan os.Signal),
		clientSet:           make(map[string]Client),
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
			s.startShutDown()
			return
		}
	}
}

func (s *Server) Stop() {
	s.signalChannel <- os.Interrupt
}

func (s *Server) startShutDown() {
	log.Info("interrupt received, starting shut down sequence")
	s.cmdHandler.Stop()
	s.disconnectClients()

}

func (s *Server) disconnectClients() {
	for _, client := range s.clientSet {
		client.Disconnect()
	}
	log.Info("all clients disconnected")
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
	client := CreateClient(conn, s.onClientDisconnected)
	s.clientSet[client.Address] = client
	go client.HandleConnection(s.cmdHandler.AppendRequest)
}

func (s *Server) onClientDisconnected(c *Client) {
	log.WithField("address", c.Address).Info("erasing client from client set")
	delete(s.clientSet, c.Address)
}
