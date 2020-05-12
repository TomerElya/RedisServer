package server

import (
	"fmt"
)

type ErrUnexpectedToken struct {
	expected, unexpected byte
}

func (e ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("Unexpected token : %b, expected : %b\n", e.unexpected, e.expected)
}

type ErrMismatchingLength struct {
	read, expected int
}

func (e ErrMismatchingLength) Error() string {
	return fmt.Sprintf("Mismatching reading length. Expected : %b, read : %b\n", e.expected, e.read)
}

type ErrArrayParsing struct{}

func (e ErrArrayParsing) Error() string {
	return fmt.Sprintf("Failed to parse command action")
}

type ErrNoCommandActionFound struct{}

func (e ErrNoCommandActionFound) Error() string {
	return fmt.Sprintf("Failed to parse command action")
}

type ErrInvalidCommandAction struct{}

func (e ErrInvalidCommandAction) Error() string {
	return fmt.Sprintf("Unable to interpret command action from given parameter")
}

type ErrUnknownMessageType struct{}

func (e ErrUnknownMessageType) Error() string {
	return fmt.Sprintf("Unrecognized message type found at beginning of value")
}

type ErrArrayLengthExtraction struct{ error }

func (e ErrArrayLengthExtraction) Error() string {
	return fmt.Sprintf("failed to extract array length. error: %v", e.error)
}

type ErrCommandNotFound struct{ command string }

func (e ErrCommandNotFound) Error() string {
	return fmt.Sprintf("command %s is not a known command", e.command)
}
