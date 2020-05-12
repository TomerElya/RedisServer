package server

import (
	"errors"
	"fmt"
)

type ErrUnexpectedToken struct {
	expected, unexpected byte
}

func (e *ErrUnexpectedToken) Error() error {
	return errors.New(fmt.Sprintf("Unexpected token : %b, expected : %b\n", e.unexpected, e.expected))
}

type ErrMismatchingLength struct {
	read, expected int
}

func (e *ErrMismatchingLength) Error() error {
	return errors.New(fmt.Sprintf("Mismatching reading length. Expected : %b, read : %b\n", e.expected, e.read))
}

type ErrArrayParsing struct{}

func (e *ErrArrayParsing) Error() error {
	return errors.New(fmt.Sprintf("Failed to parse command action"))
}

type ErrNoCommandActionFound struct{}

func (e *ErrNoCommandActionFound) Error() error {
	return errors.New(fmt.Sprintf("Failed to parse command action"))
}

type ErrInvalidCommandAction struct{}

func (e *ErrInvalidCommandAction) Error() error {
	return errors.New(fmt.Sprintf("Unable to interpret command action from given parameter"))
}

type ErrUnknownMessageType struct{}

func (e *ErrUnknownMessageType) Error() error {
	return errors.New(fmt.Sprintf("Unrecognized message type found at beginning of value"))
}

type ErrArrayLengthExtraction struct{ error }

func (e *ErrArrayLengthExtraction) Error() error {
	return errors.New(fmt.Sprintf("failed to extract array length. error: %v", e.error))
}
