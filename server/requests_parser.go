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

type RequestsParser struct {
	parserMap map[byte]func(reader *bufio.Reader) (param, error)
}

type param struct {
	messageType   byte
	value         string
	chainedParams []param
}

func createRequestParser() RequestsParser {
	ch := RequestsParser{}
	ch.initParserMap()
	return ch
}

func (ch *RequestsParser) initParserMap() {
	ch.parserMap = map[byte]func(reader *bufio.Reader) (param, error){
		str: ch.parseStr,
		num: ch.parseNum,
		blk: ch.parseBlk,
		arr: ch.parseArr,
	}
}

func (ch *RequestsParser) ConstructRequest(reader *bufio.Reader) (Request, error) {
	messageType, err := reader.ReadByte()
	if err != nil {
		return Request{}, err
	}
	parsingFunction, err := ch.getParsingFunction(messageType)
	if err != nil {
		return Request{}, err
	}
	param, err := parsingFunction(reader)
	if err != nil {
		return Request{}, err
	}
	return ch.initializeRequest(param)
}

func (ch *RequestsParser) initializeRequest(reqParam param) (Request, error) {
	switch reqParam.messageType {
	case arr:
		if err := validateArrRequest(reqParam); err != nil {
			return Request{}, err
		}
		return Request{action: reqParam.chainedParams[0].value, params: reqParam.chainedParams}, nil
	case str, blk:
		return Request{action: reqParam.value, params: nil}, nil
	case num:
		return Request{}, ErrInvalidCommandAction{}
	default:
		return Request{}, ErrUnknownMessageType{}
	}
}

func (ch *RequestsParser) parseStr(reader *bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if err := validateLF(reader); err != nil {
		return param{}, err
	}
	return param{value: string(data), messageType: str}, nil
}

func (ch *RequestsParser) parseNum(reader *bufio.Reader) (param, error) {
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

func (ch *RequestsParser) parseBlk(reader *bufio.Reader) (param, error) {
	length, err := extractNumber(reader)
	str := make([]byte, length)
	read, err := reader.Read(str)
	if err != nil {
		return param{}, err
	} else if read != length {
		return param{}, ErrMismatchingLength{read, length}
	}
	if err = validateCLRF(reader); err != nil {
		return param{}, err
	}
	return param{value: string(str), messageType: blk}, nil
}

func (ch *RequestsParser) parseArr(reader *bufio.Reader) (param, error) {
	arrSize, err := extractNumber(reader)
	if err != nil {
		return param{}, err
	}
	params := make([]param, arrSize)
	for i := 0; i < arrSize; i++ {
		messageType, err := reader.ReadByte()
		if err != nil {
			return param{}, err
		}
		parsingFunction, err := ch.getParsingFunction(messageType)
		if err != nil {
			return param{}, err
		}
		newParam, err := parsingFunction(reader)
		if err != nil {
			return param{}, err
		}
		params[i] = newParam
	}
	return param{chainedParams: params, messageType: arr}, nil
}

func extractNumber(reader *bufio.Reader) (int, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return 0, ErrArrayLengthExtraction{err}
	}
	data = data[:len(data)-1]
	length, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, ErrArrayLengthExtraction{err}
	}
	return length, validateLF(reader)
}

func (ch *RequestsParser) getParsingFunction(messageType byte) (func(reader *bufio.Reader) (param, error), error) {
	if function, ok := ch.parserMap[messageType]; !ok {
		return nil, ErrUnknownMessageType{}
	} else {
		return function, nil
	}
}
