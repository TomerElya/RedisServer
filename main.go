package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tomer/RedisServer/server"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	s := server.CreateServer("localhost", 3000)
	s.StartAndListen()

}
