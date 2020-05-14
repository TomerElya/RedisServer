package server

import log "github.com/sirupsen/logrus"

type Store struct {
	store            map[string]Param
	actionMap        map[string]func(request StoreRequest)
	IncomingRequests chan StoreRequest
	StopChan         chan interface{}
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
	store := Store{
		store:            map[string]Param{},
		IncomingRequests: make(chan StoreRequest),
		StopChan:         make(chan interface{}),
	}
	store.initializeActionMap()
	return store
}

func (s *Store) initializeActionMap() {
	s.actionMap = map[string]func(request StoreRequest){
		"get": s.Get,
		"set": s.Set,
	}
}

func (s *Store) Start() {
	go s.listen()
}

func (s *Store) Stop() {
	s.StopChan <- 0
}

func (s *Store) listen() {
	for {
		select {
		case req := <-s.IncomingRequests:
			s.actionMap[req.action](req)
		case <-s.StopChan:
			s.handleInterrupt()
			return
		}
	}
}

func (s *Store) handleInterrupt() {
	log.Info("store received interrupt")
	for req := range s.IncomingRequests {
		s.actionMap[req.action](req)
	}
	log.Info("all left requests are handled, closing...")
}

func (s *Store) Exists(command string) bool {
	_, ok := s.actionMap[command]
	return ok
}

func (s *Store) Get(request StoreRequest) {
	val, ok := s.store[request.params[1].value]
	if !ok {
		request.responseChan <- StoreResponse{Param{}, ErrKeyNotFound{key: request.params[1].value}}
	}
	request.responseChan <- StoreResponse{val, nil}
}

func (s *Store) Set(request StoreRequest) {
	s.store[request.params[1].value] = request.params[2]
	request.responseChan <- StoreResponse{Param{value: "OK", messageType: str}, nil}
}
