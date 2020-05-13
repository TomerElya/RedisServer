package server

type Store struct {
	store            map[string]Param
	IncomingRequests chan StoreRequest
}

type StoreResponse struct {
	Param
	error
}

type StoreRequest struct {
	Param
	responseChan chan StoreResponse
}

func CreateStore() Store {
	return Store{
		store:            map[string]Param{},
		IncomingRequests: make(chan StoreRequest),
	}
}

func (s *Store) Start() {
	go s.listen()
}

func (s *Store) listen() {
	select {
	case req := <-s.incomingRequests:

	}
}

func (s *Store) Get(request StoreRequest) {
	val, ok := s.store[request.Param.chainedParams[1].value]
	if !ok {
		request.responseChan <- StoreResponse{Param{}, ErrKeyNotFound{}}
	}
	request.responseChan <- StoreResponse{val, nil}
}
