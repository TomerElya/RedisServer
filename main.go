package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tomer/RedisServer/server"
	"io"
	"os"
	"runtime/pprof"
)

func main() {
	err := pprof.StartCPUProfile(createProfilingFile())
	if err != nil {
		panic(err)
	}
	log.SetOutput(os.Stdout)
	s := server.CreateServer("localhost", 3000)
	s.StartAndListen()
	pprof.StopCPUProfile()
}

func createProfilingFile() io.Writer {
	file, err := os.Create("./profile.pprof")
	if err != nil {
		panic(err)
	}
	return file
}
