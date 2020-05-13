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
	isConnected int32
}

func CreateClient(conn net.Conn) Client {
	client := Client{
		reader:      bufio.NewReader(conn),
		mutex:       sync.Mutex{},
		logger:      log.WithField("address", conn.RemoteAddr().String()),
		isConnected: 0,
		conn:        conn,
	}
	client.logger.Info("new client created")
	return client
}

func (c *Client) Disconnect(err error) {
	if err == io.EOF {
		c.logger.Info("client disconnected")
		atomic.StoreInt32(&c.isConnected, 1)
	} else {
		c.logger.WithError(err).Error("error received while listening to connection. Disconnecting...")
		_, err = c.conn.Write([]byte(err.Error()))
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
	param := Param{value: err.Error(), messageType: err1}
	err = c.write(param)
	if err != nil {
		c.logger.WithError(err).WithField("error", err).Error("failed to write error to client")
	}
}

func (c *Client) write(param Param) error {
	if atomic.LoadInt32(&c.isConnected) != 0 {
		return ErrConnectionClosedWrite{}
	}
	c.mutex.Lock()
	written, err := c.conn.Write(param.ToBytes())
	c.mutex.Unlock()
	if len(param.value) != written {
		return ErrIncompleteWrite{written: written, expected: len(param.value)}
	}
	return err
}
