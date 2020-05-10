package server

import (
	"bufio"
)

type request struct {
	action string
	params []string
}

func constructRequest(reader *bufio.Reader) (request, error) {
	req := request{}
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		return req, err
	}

}
