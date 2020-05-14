package server

import log "github.com/sirupsen/logrus"

type CommandHandler struct {
	stopChan chan bool
	store    Store
}

func CreateCommandHandler() CommandHandler {
	cmdHandler := CommandHandler{
		stopChan: make(chan bool),
		store:    CreateStore(),
	}
	cmdHandler.store.Start()
	return cmdHandler
}

func (ch *CommandHandler) AppendRequest(req Request) {
	log.WithField("request action", req.action).WithField("address", req.client.conn.RemoteAddr().String()).
		Info("new request received")
	ok := ch.store.Exists(req.action)
	if !ok {
		req.client.WriteError(ErrCommandNotFound{command: req.action})
	} else {
		go ch.handleRequest(req)
	}
}

func (ch *CommandHandler) Stop() {
	log.Info("command handler received shut down interrupt")
	ch.store.Stop()
}

func (ch *CommandHandler) handleRequest(req Request) {
	storeRequest := StoreRequest{responseChan: make(chan StoreResponse), Request: req}
	ch.store.IncomingRequests <- storeRequest
	storeResponse := <-storeRequest.responseChan
	if storeResponse.error != nil {
		req.client.WriteError(storeResponse.error)
	} else {
		req.client.WriteResponse(storeResponse.Param)
	}
}
