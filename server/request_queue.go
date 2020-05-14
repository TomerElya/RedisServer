package server

type RequestQueue struct {
	requests chan Request
}

var QueueBufferSize = 10000

func CreateRequestQueue() RequestQueue {
	rq := RequestQueue{
		requests: make(chan Request, QueueBufferSize),
	}
	return rq
}

func (rq *RequestQueue) Push(request Request) error {
	if len(rq.requests) >= QueueBufferSize {
		return ErrQueueFull{}
	}
	rq.requests <- request
	return nil
}

func (rq *RequestQueue) Pop() Request {
	return <-rq.requests
}
