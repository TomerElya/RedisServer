package server

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	conn        net.Conn
	reader      *bufio.Reader
	logger      *log.Entry
	mutex       sync.Mutex
	stopChan    chan interface{}
	reqChan     chan Request
	isConnected int32
}

func CreateClient(conn net.Conn) Client {
	client := Client{
		reader:      bufio.NewReader(conn),
		mutex:       sync.Mutex{},
		logger:      log.WithField("address", conn.RemoteAddr().String()),
		isConnected: 0,
		conn:        conn,
		stopChan:    make(chan interface{}),
		reqChan:     make(chan Request),
	}
	client.logger.Info("new client created")
	return client
}

func (c *Client) HandleConnection(appendRequest func(request Request)) {
	for {
		select {
		case req := <-c.reqChan:
			appendRequest(req)
		case <-c.stopChan:
			atomic.StoreInt32(&c.isConnected, 1)
			return
		}
	}
}

func (c *Client) processRequests() {
	var err error = nil
	var req Request
	for err == nil && atomic.LoadInt32(&c.isConnected) == 0 {
		req, err = ConstructRequest(c.reader)
		if err == nil {
			req.client = c
			c.reqChan <- req
		}
	}
}

func (c *Client) Disconnect(err error) {
	if err == io.EOF {
		c.logger.Info("client disconnected")
		atomic.StoreInt32(&c.isConnected, 1)
	} else {
		c.logger.WithError(err).Error("error received while listening to connection. Disconnecting...")
		err = c.write(Param{value: err.Error(), messageType: err1})
		if err != nil {
			c.logger.WithError(err).Error("failed to write error to client while disconnecting")
		}
		err = c.conn.Close()
		if err != nil {
			c.logger.WithError(err).Error("error while trying to close connection")
		}
	}
}

func (c *Client) WriteError(err error) {
	param := Param{messageType: err1, value: err.Error()}
	c.logger.WithError(err).Info("writing error to client")
	writeErr := c.write(param)
	if err != nil {
		c.logger.WithError(writeErr).WithField("error", err).Error("failed to write error to client")
	}
}

func (c *Client) WriteResponse(param Param) {
	c.logger.WithField("response", param.value).WithField("message type", param.messageType).
		Info("writing successful response to client")
	err := c.write(param)
	if err != nil {
		c.logger.WithError(err).Error("failed to write response to client")
	}
}

func (c *Client) write(param Param) error {
	if atomic.LoadInt32(&c.isConnected) != 0 {
		return ErrConnectionClosedWrite{}
	}
	paramBytes := param.ToBytes()
	c.mutex.Lock()
	written, err := c.conn.Write(paramBytes)
	c.mutex.Unlock()
	if len(paramBytes) != written {
		return ErrIncompleteWrite{written: written, expected: len(paramBytes)}
	}
	return err
}
