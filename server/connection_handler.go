package server

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"net"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	logger logrus.Logger
}

func CreateClient(conn net.Conn) Client {
	client := Client{}
	client.reader = bufio.NewReader(conn)
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
