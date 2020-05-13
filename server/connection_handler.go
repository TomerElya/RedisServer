package server

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	logger *log.Entry
	mutex  sync.Mutex
}

func CreateClient(conn net.Conn) Client {
	client := Client{
		reader: bufio.NewReader(conn),
		mutex:  sync.Mutex{},
		logger: log.WithField("address", conn.RemoteAddr().String()),
	}
	return client
}

func (c *Client) DisconnectWithError(err error) {
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

func (c *Client) WriteError(err error) {
	c.logger.Info("new connection established")
	_, err = c.conn.Write([]byte(err.Error()))
	if err != nil {
		c.logger.WithError(err).WithField("error", err).Error("failed to write error to client")
	}
}
