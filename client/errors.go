package client

import (
	"errors"
	"fmt"
)

// Common errors that can be returned by the client.
var (
	// ErrEmptyInput is returned when the input text is empty.
	ErrEmptyInput = errors.New("input text is empty")

	// ErrNilMessage is returned when the protobuf message is nil.
	ErrNilMessage = errors.New("protobuf message cannot be nil")

	// ErrNoAnnotators is returned when no annotators are specified.
	ErrNoAnnotators = errors.New("at least one annotator is required")
)

// AnnotatorError represents an error related to annotator configuration.
type AnnotatorError struct {
	Annotator string
	Message   string
}

func (e *AnnotatorError) Error() string {
	return fmt.Sprintf("annotator error [%s]: %s", e.Annotator, e.Message)
}

// ServerError represents an error from the CoreNLP server.
type ServerError struct {
	URL        string
	StatusCode int
	Message    string
}

func (e *ServerError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("server error [%s] (status %d): %s", e.URL, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("server error [%s]: %s", e.URL, e.Message)
}

// ParseError represents an error during protobuf parsing.
type ParseError struct {
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("parse error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// CommandError represents an error from executing the CoreNLP command.
type CommandError struct {
	Command string
	Stderr  string
	Err     error
}

func (e *CommandError) Error() string {
	if e.Stderr != "" {
		return fmt.Sprintf("command error [%s]: %v\nstderr: %s", e.Command, e.Err, e.Stderr)
	}
	return fmt.Sprintf("command error [%s]: %v", e.Command, e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}
