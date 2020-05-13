package server

type Store struct {
	store            map[string]Param
	actionMap        map[string]func(request StoreRequest)
	IncomingRequests chan StoreRequest
}

type StoreResponse struct {
	Param
	error
}

type StoreRequest struct {
	Request
	responseChan chan StoreResponse
}

func CreateStore() Store {
	return Store{
		store:            map[string]Param{},
		actionMap:        map[string]func(request StoreRequest){},
		IncomingRequests: make(chan StoreRequest),
	}
}

func (s *Store) initializeActionMap() {
	s.actionMap = map[string]func(request StoreRequest){
		"get": s.Get,
	}
}

func (s *Store) Start() {
	go s.listen()
}

func (s *Store) listen() {
	select {
	case req := <-s.IncomingRequests:
		s.actionMap[req.action](req)
	}
}

func (s *Store) Get(request StoreRequest) {
	val, ok := s.store[request.params[1].value]
	if !ok {
		request.responseChan <- StoreResponse{Param{}, ErrKeyNotFound{}}
	}
	request.responseChan <- StoreResponse{val, nil}
}
