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
		blk: parseBlk,
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
	if err := verifyLF(reader); err != nil {
		return param{}, err
	}
	return param{value: string(data), messageType: str}, nil
}

func parseNum(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if err := verifyLF(reader); err != nil {
		return param{}, err
	}
	number, err := strconv.Atoi(string(data))
	if err != nil {
		return param{}, err
	}
	return param{value: number, messageType: str}, nil
}

func parseBlk(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	length, err := strconv.Atoi(string(data))
	if err != nil {
		return param{}, err
	}
	if err := verifyLF(reader); err != nil {
		return param{}, err
	}
	str := make([]byte, length)
	read, err := reader.Read(str)
	if err != nil {
		return param{}, err
	} else if read != length {
		return param{}, MismatchingLength{read, length}.Error()
	}
	if err = verifyCLRF(reader); err != nil {
		return param{}, err
	}
	return param{value: str, messageType: blk}, nil
}

func verifyLF(reader bufio.Reader) error {
	if nextByte, err := reader.ReadByte(); err != nil {
		return err
	} else if nextByte != '\n' {
		return UnexpectedToken{'\n', nextByte}.Error()
	}
	return nil
}

func verifyCR(reader bufio.Reader) error {
	if nextByte, err := reader.ReadByte(); err != nil {
		return err
	} else if nextByte != '\r' {
		return UnexpectedToken{'\r', nextByte}.Error()
	}
	return nil
}

func verifyCLRF(reader bufio.Reader) error {
	err := verifyCR(reader)
	if err != nil {
		err = verifyLF(reader)
	}
	return err
}
