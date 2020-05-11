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

type InvalidCommandActionError struct{}

func (e *InvalidCommandActionError) Error() error {
	return errors.New(fmt.Sprintf("Unable to interpret command action from given parameter"))
}

type UnknownMessageTypeError struct{}

func (e *UnknownMessageTypeError) Error() error {
	return errors.New(fmt.Sprintf("Unrecognized message type found at beginning of value"))
}

type ArrayLengthExtractionError struct{ error }

func (e *ArrayLengthExtractionError) Error() error {
	return errors.New(fmt.Sprintf("failed to extract array length. error: %v", e.error))
}
