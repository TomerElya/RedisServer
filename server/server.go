package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	address string
	port    int
	mux     *http.ServeMux
}

func CreateServer(address string, port int) Server {
	return Server{address: address, port: port}
}

func (s *Server) Start() bool {
	address := fmt.Sprintf("%s:%d", s.address, s.port)
	s.createMuxHandler()
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.WithError(err).WithField("address", address)
		return false
	}
	log.WithField("address", address).Info("http server successfully started")
	return true
}

func (s *Server) createMuxHandler() {
	s.mux = http.NewServeMux()

}
