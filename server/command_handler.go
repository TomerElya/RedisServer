package server

type CommandHandler struct {
	stopChan   chan bool
	store      Store
	commandMap map[string]func(req Request)
}

func CreateCommandHandler() CommandHandler {
	cmdHandler := CommandHandler{
		stopChan: make(chan bool),
		store:    CreateStore(),
	}
	cmdHandler.initializeCommandMap()
	cmdHandler.store.Start()
	return cmdHandler
}

func (ch *CommandHandler) initializeCommandMap() {
	ch.commandMap = map[string]func(req Request){
		"get": ch.handleGet,
	}
}

func (ch *CommandHandler) AppendRequest(req Request) {
	handlerFunc, ok := ch.commandMap[req.action]
	if !ok {
		req.client.WriteError(ErrCommandNotFound{command: req.action})
	} else {
		handlerFunc(req)
	}
}

func (ch *CommandHandler) handleGet(req Request) {
	storeRequest := StoreRequest{responseChan: make(chan StoreResponse), Request: req}
	ch.store.IncomingRequests <- storeRequest
	storeResponse := <-storeRequest.responseChan
	if storeResponse.error != nil {
		req.client.WriteError(storeResponse.error)
	} else {
		req.client.WriteResponse(storeResponse.Param)
	}
}
