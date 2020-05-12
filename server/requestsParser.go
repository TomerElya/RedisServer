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

type reqParser struct {
	parserMap map[byte]func(reader *bufio.Reader) (param, error)
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

func createRequestParser() reqParser {
	ch := reqParser{}
	ch.initParserMap()
	return ch
}

func (ch *reqParser) initParserMap() {
	ch.parserMap = map[byte]func(reader *bufio.Reader) (param, error){
		str: ch.parseStr,
		num: ch.parseNum,
		blk: ch.parseBlk,
		arr: ch.parseArr,
	}
}

func (ch *reqParser) constructRequest(reader *bufio.Reader) (request, error) {
	messageType, err := reader.ReadByte()
	if err != nil {
		return request{}, err
	}
	parsingFunction, err := ch.getParsingFunction(messageType)
	if err != nil {
		return request{}, err
	}
	param, err := parsingFunction(reader)
	if err != nil {
		return request{}, err
	}
	return ch.initializeRequest(param)
}

func (ch *reqParser) initializeRequest(reqParam param) (request, error) {
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

func (ch *reqParser) parseStr(reader *bufio.Reader) (param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return param{}, err
	}
	if err := validateLF(reader); err != nil {
		return param{}, err
	}
	return param{value: string(data), messageType: str}, nil
}

func (ch *reqParser) parseNum(reader *bufio.Reader) (param, error) {
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

func (ch *reqParser) parseBlk(reader *bufio.Reader) (param, error) {
	length, err := extractNumber(reader)
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

func (ch *reqParser) parseArr(reader *bufio.Reader) (param, error) {
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
		return 0, (&ArrayLengthExtractionError{err}).Error()
	}
	data = data[:len(data)-1]
	length, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, (&ArrayLengthExtractionError{err}).Error()
	}
	return length, validateLF(reader)
}

func (ch *reqParser) getParsingFunction(messageType byte) (func(reader *bufio.Reader) (param, error), error) {
	if function, ok := ch.parserMap[messageType]; !ok {
		return nil, (&UnknownMessageTypeError{}).Error()
	} else {
		return function, nil
	}
}
