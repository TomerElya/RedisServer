package server

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
	ok := ch.store.Exists(req.action)
	if !ok {
		req.client.WriteError(ErrCommandNotFound{command: req.action})
	} else {
		ch.handleRequest(req)
	}
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
