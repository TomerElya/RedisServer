package server

type CommandHandler struct {
	incomingRequests chan commandForm
	stopChan         chan bool
	store
	commandMap map[string]func(req Request)
}

type commandForm struct {
	commandFunc  func(req Request)
	request      Request
	responseChan chan string
}

func CreateCommandHandler() CommandHandler {
	cmdHandler := CommandHandler{incomingRequests: make(chan commandForm)}
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

func (ch *CommandHandler) AppendRequest(req Request) {
	handlerFunc, ok := ch.commandMap[req.action]
	if !ok {
		req.client.WriteError(ErrCommandNotFound{command: req.action})
	}
	responseChan := make(chan string)
	cmdForm := commandForm{commandFunc: handlerFunc, request: req, responseChan: responseChan}
	ch.incomingRequests <- cmdForm
	response := <-cmdForm.responseChan
	req.client.write([]byte(response))
}

func (ch *CommandHandler) process() {
	select {
	case cmdForm := <-ch.incomingRequests:
		cmdForm.commandFunc(cmdForm.request)
	case <-ch.stopChan:
		break
	}
}

func (ch *CommandHandler) handleGet(req Request) {
	value, ok := ch.store
}
