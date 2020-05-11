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
	parserMap map[byte]func(reader bufio.Reader) ([]param, error)
}

type request struct {
	action string
	params []param
}

type param struct {
	messageType byte
	value       interface{}
}

func CreateCommandHandler() cmdHandler {
	ch := cmdHandler{}
	ch.initParserMap()
}

func (ch *cmdHandler) initParserMap() map[byte]func(reader bufio.Reader) ([]param, error) {
	return map[byte]func(reader bufio.Reader) ([]param, error){
		str: ch.parseStr,
		num: ch.parseNum,
		blk: ch.parseBlk,
		arr: ch.parseArr,
	}
}

func (ch *cmdHandler) constructRequest(reader bufio.Reader) (request, error) {
	req := request{}
	messageType, err := reader.ReadByte()
	if err != nil {
		return req, err
	}
}

func (ch *cmdHandler) parseStr(reader bufio.Reader) ([]param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return []param{}, err
	}
	if err := verifyLF(reader); err != nil {
		return []param{}, err
	}
	return []param{{value: string(data), messageType: str}}, nil
}

func (ch *cmdHandler) parseNum(reader bufio.Reader) ([]param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return []param{}, err
	}
	if err := verifyLF(reader); err != nil {
		return []param{}, err
	}
	number, err := strconv.Atoi(string(data))
	if err != nil {
		return []param{}, err
	}
	return []param{{value: number, messageType: num}}, nil
}

func (ch *cmdHandler) parseBlk(reader bufio.Reader) ([]param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return []param{}, err
	}
	length, err := strconv.Atoi(string(data))
	if err != nil {
		return []param{}, err
	}
	if err := verifyLF(reader); err != nil {
		return []param{}, err
	}
	str := make([]byte, length)
	read, err := reader.Read(str)
	if err != nil {
		return []param{}, err
	} else if read != length {
		return []param{}, MismatchingLength{read, length}.Error()
	}
	if err = verifyCLRF(reader); err != nil {
		return []param{}, err
	}
	return []param{{value: str, messageType: blk}}, nil
}

func (ch *cmdHandler) parseArr(reader bufio.Reader) ([]param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return []param{}, err
	}
	arrSize, err := strconv.Atoi(string(data))
	if err != nil {
		return []param{}, err
	}
	params := make([]param, arrSize)
	if err := verifyLF(reader); err != nil {
		return []param{}, err
	}
	for i := 0; i < arrSize; i++ {
		messageType, err := reader.ReadByte()
		if err != nil {
			return []param{}, err
		}
		newParams, err := ch.parserMap[messageType](reader)
		if err != nil {
			return []param{}, err
		}
		params = append(params, newParams...)
	}
	return params, nil
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
