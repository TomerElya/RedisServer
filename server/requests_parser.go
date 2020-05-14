package server

import (
	"bufio"
	"strconv"
	"strings"
)

var parserMap = map[byte]func(reader *bufio.Reader) (Param, error){
	str: parseStr,
	num: parseNum,
	blk: parseBlk,
	arr: parseArr,
}

func ConstructRequest(reader *bufio.Reader) (Request, error) {
	messageType, err := reader.ReadByte()
	if err != nil {
		return Request{}, err
	}
	parsingFunction, err := getParsingFunction(messageType)
	if err != nil {
		return Request{}, err
	}
	Param, err := parsingFunction(reader)
	if err != nil {
		return Request{}, err
	}
	return initializeRequest(Param)
}

func initializeRequest(reqParam Param) (Request, error) {
	var resultedRequest Request
	switch reqParam.messageType {
	case arr:
		if err := validateArrRequest(reqParam); err != nil {
			return Request{}, err
		}
		resultedRequest = Request{action: reqParam.chainedParams[0].value, params: reqParam.chainedParams}
	case str, blk:
		resultedRequest = Request{action: reqParam.value, params: nil}
	case num:
		return Request{}, ErrInvalidCommandAction{}
	default:
		return Request{}, ErrUnknownMessageType{}
	}
	resultedRequest.action = strings.ToLower(resultedRequest.action)
	return resultedRequest, nil
}

func parseStr(reader *bufio.Reader) (Param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return Param{}, err
	}
	if err := validateLF(reader); err != nil {
		return Param{}, err
	}
	return Param{value: string(data), messageType: str}, nil
}

func parseNum(reader *bufio.Reader) (Param, error) {
	data, err := reader.ReadBytes('\r')
	if err != nil {
		return Param{}, err
	}
	if err := validateLF(reader); err != nil {
		return Param{}, err
	}
	stringData := string(data)
	_, err = strconv.Atoi(stringData)
	if err != nil {
		return Param{}, err
	}
	return Param{value: stringData, messageType: num}, nil
}

func parseBlk(reader *bufio.Reader) (Param, error) {
	length, err := extractNumber(reader)
	str := make([]byte, length)
	read, err := reader.Read(str)
	if err != nil {
		return Param{}, err
	} else if read != length {
		return Param{}, ErrMismatchingLength{read, length}
	}
	if err = validateCLRF(reader); err != nil {
		return Param{}, err
	}
	return Param{value: string(str), messageType: blk}, nil
}

func parseArr(reader *bufio.Reader) (Param, error) {
	arrSize, err := extractNumber(reader)
	if err != nil {
		return Param{}, err
	}
	Params := make([]Param, arrSize)
	for i := 0; i < arrSize; i++ {
		messageType, err := reader.ReadByte()
		if err != nil {
			return Param{}, err
		}
		parsingFunction, err := getParsingFunction(messageType)
		if err != nil {
			return Param{}, err
		}
		newParam, err := parsingFunction(reader)
		if err != nil {
			return Param{}, err
		}
		Params[i] = newParam
	}
	return Param{chainedParams: Params, messageType: arr}, nil
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

func getParsingFunction(messageType byte) (func(reader *bufio.Reader) (Param, error), error) {
	if function, ok := parserMap[messageType]; !ok {
		return nil, ErrUnknownMessageType{}
	} else {
		return function, nil
	}
}
