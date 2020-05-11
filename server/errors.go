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

type ArrayParsingError struct{}

func (e *ArrayParsingError) Error() error {
	return errors.New(fmt.Sprintf("Failed to parse command action"))
}

type NoCommandActionFoundError struct{}

func (e *NoCommandActionFoundError) Error() error {
	return errors.New(fmt.Sprintf("Failed to parse command action"))
}
