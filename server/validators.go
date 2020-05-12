package server

import "bufio"

func validateArrRequest(reqParam param) error {
	if len(reqParam.chainedParams) < 1 {
		return ErrArrayParsing{}
	}
	if reqParam.chainedParams[0].messageType != str && reqParam.chainedParams[0].messageType != blk {
		return ErrNoCommandActionFound{}
	}
	return nil
}

func validateLF(reader *bufio.Reader) error {
	if nextByte, err := reader.ReadByte(); err != nil {
		return err
	} else if nextByte != '\n' {
		return ErrUnexpectedToken{'\n', nextByte}
	}
	return nil
}

func validateCR(reader *bufio.Reader) error {
	if nextByte, err := reader.ReadByte(); err != nil {
		return err
	} else if nextByte != '\r' {
		return ErrUnexpectedToken{'\r', nextByte}
	}
	return nil
}

func validateCLRF(reader *bufio.Reader) error {
	err := validateCR(reader)
	if err == nil {
		err = validateLF(reader)
	}
	return err
}
