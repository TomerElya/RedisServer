package server

import "bufio"

func validateArrRequest(reqParam param) error {
	if len(reqParam.chainedParams) < 1 {
		return (&ArrayParsingError{}).Error()
	}
	if reqParam.chainedParams[0].messageType != str && reqParam.chainedParams[0].messageType != blk {
		return (&NoCommandActionFoundError{}).Error()
	}
	return nil
}

func validateLF(reader *bufio.Reader) error {
	if nextByte, err := reader.ReadByte(); err != nil {
		return err
	} else if nextByte != '\n' {
		return (&UnexpectedToken{'\n', nextByte}).Error()
	}
	return nil
}

func validateCR(reader *bufio.Reader) error {
	if nextByte, err := reader.ReadByte(); err != nil {
		return err
	} else if nextByte != '\r' {
		return (&UnexpectedToken{'\r', nextByte}).Error()
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
