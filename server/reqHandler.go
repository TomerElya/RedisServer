package server

import (
	"bufio"
	"strconv"
)

const (
	str = '+'
	num = ':'
	blk = '$'
	arr = '*'
)

type cmdHandler struct {
	parserMap map[byte]func(reader bufio.Reader) (param, error)
}

type request struct {
	action string
	params []string
}

type param struct {
	messageType byte
	value       interface{}
}

func CreateCommandHandler() cmdHandler {
	return cmdHandler{parserMap: initParserMap()}
}

func initParserMap() map[byte]func(reader bufio.Reader) (param, error) {
	return map[byte]func(reader bufio.Reader) (param, error){
		str: parseStr,
		num: parseNum,
	}
}

func (ch *cmdHandler) constructRequest(reader bufio.Reader) (request, error) {
	req := request{}
	messageType, err := reader.ReadByte()
	if err != nil {
		return req, err
	}

}

func parseStr(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if nextByte, err := reader.ReadByte(); err != nil {
		return param{}, err
	} else if nextByte != '\n' {
		return param{}, UnexpectedToken{'\n', nextByte}.Error()
	}
	return param{value: string(data), messageType: str}, nil
}

func parseNum(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if nextByte, err := reader.ReadByte(); err != nil {
		return param{}, err
	} else if nextByte != '\n' {
		return param{}, UnexpectedToken{'\n', nextByte}.Error()
	}
	number, err := strconv.Atoi(string(data))
	if err != nil {
		return param{}, err
	}
	return param{value: number, messageType: str}, nil
}
