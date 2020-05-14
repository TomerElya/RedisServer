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
	Address          string
	conn             net.Conn
	reader           *bufio.Reader
	logger           *log.Entry
	mutex            sync.Mutex
	stopChan         chan interface{}
	reqChan          chan Request
	isConnected      int32
	notifyDisconnect func(*Client)
}

func CreateClient(conn net.Conn, notifyDisconnect func(*Client)) Client {
	client := Client{
		Address:          conn.RemoteAddr().String(),
		reader:           bufio.NewReader(conn),
		mutex:            sync.Mutex{},
		logger:           log.WithField("address", conn.RemoteAddr().String()),
		isConnected:      0,
		conn:             conn,
		stopChan:         make(chan interface{}),
		reqChan:          make(chan Request),
		notifyDisconnect: notifyDisconnect,
	}
	client.logger.Info("new client connected")
	return client
}

func (c *Client) HandleConnection(appendRequest func(request Request)) {
	go c.processRequests()
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
	if err != nil {
		c.DisconnectWithError(err)
	}
}

func (c *Client) Disconnect() {
	c.logger.Info("disconnecting client")
	err := c.conn.Close()
	if err != nil {
		c.logger.WithError(err).Error("error while trying to close connection")
	}
	c.notifyDisconnect(c)
}

func (c *Client) DisconnectWithError(err error) {
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
	c.notifyDisconnect(c)
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
