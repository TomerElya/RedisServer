package server

type CommandHandler struct {
	incomingRequests chan commandForm
	stopChan         chan bool
	store            Store
	commandMap       map[string]func(req Request)
}

type commandForm struct {
	commandFunc  func(req Request)
	request      Request
	responseChan chan string
}

func CreateCommandHandler() CommandHandler {
	cmdHandler := CommandHandler{
		incomingRequests: make(chan commandForm),
		stopChan:         make(chan bool), store: CreateStore(),
	}
	cmdHandler.initializeCommandMap()
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
		return
	}
	responseChan := make(chan string)
	cmdForm := commandForm{commandFunc: handlerFunc, request: req, responseChan: responseChan}
	ch.incomingRequests <- cmdForm
	response := <-cmdForm.responseChan
	req.client.write([]byte(response))
}

func (ch *CommandHandler) handleGet(req Request) {
	value, err := ch.store.Get(req.params[1].value)
}
