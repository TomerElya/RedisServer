package server

import (
	"errors"
	"fmt"
)

type UnexpectedToken struct {
	expected, unexpected byte
}

func (e *UnexpectedToken) Error() error {
	return errors.New(fmt.Sprintf("Unexpected token : %b, expected : %b\n", e.unexpected, e.expected))
}

type MismatchingLength struct {
	read, expected int
}

func (e *MismatchingLength) Error() error {
	return errors.New(fmt.Sprintf("Mismatching reading length. Expected : %b, read : %b\n", e.expected, e.read))
}
