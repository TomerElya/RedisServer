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
	params []param
}

type param struct {
	messageType   byte
	value         string
	chainedParams []param
}

func CreateCommandHandler() cmdHandler {
	ch := cmdHandler{}
	ch.initParserMap()
	return ch
}

func (ch *cmdHandler) initParserMap() map[byte]func(reader bufio.Reader) (param, error) {
	return map[byte]func(reader bufio.Reader) (param, error){
		str: ch.parseStr,
		num: ch.parseNum,
		blk: ch.parseBlk,
		arr: ch.parseArr,
	}
}

func (ch *cmdHandler) constructRequest(reader bufio.Reader) (request, error) {
	messageType, err := reader.ReadByte()
	if err != nil {
		return request{}, err
	}
	param, err := ch.parserMap[messageType](reader)
	if err != nil {
		return request{}, err
	}
	return ch.initializeRequest(param)
}

func (ch *cmdHandler) initializeRequest(reqParam param) (request, error) {
	switch reqParam.messageType {
	case arr:
		if err := validateArrRequest(reqParam); err != nil {
			return request{}, err
		}
		return request{action: reqParam.chainedParams[0].value, params: reqParam.chainedParams}, nil
	case str, blk:
		return request{action: reqParam.value, params: nil}, nil
	case num:
		return request{}, (&InvalidCommandActionError{}).Error()
	default:
		return request{}, (&UnknownMessageTypeError{}).Error()
	}
}

func (ch *cmdHandler) parseStr(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if err := validateLF(reader); err != nil {
		return param{}, err
	}
	return param{value: string(data), messageType: str}, nil
}

func (ch *cmdHandler) parseNum(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if err := validateLF(reader); err != nil {
		return param{}, err
	}
	stringData := string(data)
	_, err = strconv.Atoi(stringData)
	if err != nil {
		return param{}, err
	}
	return param{value: stringData, messageType: num}, nil
}

func (ch *cmdHandler) parseBlk(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	length, err := strconv.Atoi(string(data))
	if err != nil {
		return param{}, err
	}
	if err := validateLF(reader); err != nil {
		return param{}, err
	}
	str := make([]byte, length)
	read, err := reader.Read(str)
	if err != nil {
		return param{}, err
	} else if read != length {
		return param{}, (&MismatchingLength{read, length}).Error()
	}
	if err = validateCLRF(reader); err != nil {
		return param{}, err
	}
	return param{value: string(str), messageType: blk}, nil
}

func (ch *cmdHandler) parseArr(reader bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	arrSize, err := strconv.Atoi(string(data))
	if err != nil {
		return param{}, err
	}
	params := make([]param, arrSize)
	if err := validateLF(reader); err != nil {
		return param{}, err
	}
	for i := 0; i < arrSize; i++ {
		messageType, err := reader.ReadByte()
		if err != nil {
			return param{}, err
		}
		newParam, err := ch.parserMap[messageType](reader)
		if err != nil {
			return param{}, err
		}
		params[i] = newParam
	}
	return param{chainedParams: params, messageType: arr}, nil
}
