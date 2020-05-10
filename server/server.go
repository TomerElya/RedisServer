package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	address string
	port    int
}

func (s *Server) Start() bool {
	address := fmt.Sprintf("%s:%d", s.address, s.port)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.WithError(err).WithField("address", address)
		return false
	}
	return true
}
