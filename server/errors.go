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
