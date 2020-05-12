package server

type CommandHandler struct {
	incomingRequests chan Request
	stopChan         chan bool
	store
	commandMap map[string]func(req Request)
}

func CreateCommandHandler() CommandHandler {
	cmdHandler := CommandHandler{incomingRequests: make(chan Request)}
	return cmdHandler
}

func (ch *CommandHandler) initializeCommandMap() {
	ch.commandMap = map[string]func(req Request){
		"get": ch.handleGet,
	}
}

func (ch *CommandHandler) Start() {
	go ch.process()
}

func (ch *CommandHandler) AppendRequest(req Request) error {
	handlerFunc, ok := ch.commandMap[req.action]
	if !ok {
		return ErrCommandNotFound{command: req.action}
	}

}

func (ch *CommandHandler) process() {
	select {
	case req := <-ch.incomingRequests:
		ch.processIncomingRequest(req)
	case <-ch.stopChan:
		break
	}
}

func (ch *CommandHandler) processIncomingRequest(req Request) {

}

func (ch *CommandHandler) handleGet(req Request) {

}
