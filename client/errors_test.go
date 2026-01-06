package client

import (
	"errors"
	"testing"
)

func TestErrorTypes(t *testing.T) {
	t.Run("AnnotatorError", func(t *testing.T) {
		err := &AnnotatorError{
			Annotator: "invalid",
			Message:   "not found",
		}
		expected := "annotator error [invalid]: not found"
		if got := err.Error(); got != expected {
			t.Errorf("AnnotatorError.Error() = %v, want %v", got, expected)
		}
	})

	t.Run("ServerError with status code", func(t *testing.T) {
		err := &ServerError{
			URL:        "http://localhost:9000",
			StatusCode: 404,
			Message:    "not found",
		}
		expected := "server error [http://localhost:9000] (status 404): not found"
		if got := err.Error(); got != expected {
			t.Errorf("ServerError.Error() = %v, want %v", got, expected)
		}
	})

	t.Run("ServerError without status code", func(t *testing.T) {
		err := &ServerError{
			URL:     "http://localhost:9000",
			Message: "connection failed",
		}
		expected := "server error [http://localhost:9000]: connection failed"
		if got := err.Error(); got != expected {
			t.Errorf("ServerError.Error() = %v, want %v", got, expected)
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		baseErr := errors.New("base error")
		err := &ParseError{
			Message: "invalid protobuf",
			Err:     baseErr,
		}
		expected := "parse error: invalid protobuf: base error"
		if got := err.Error(); got != expected {
			t.Errorf("ParseError.Error() = %v, want %v", got, expected)
		}

		// Test Unwrap
		if got := err.Unwrap(); got != baseErr {
			t.Errorf("ParseError.Unwrap() = %v, want %v", got, baseErr)
		}
	})

	t.Run("CommandError", func(t *testing.T) {
		baseErr := errors.New("exit code 1")
		err := &CommandError{
			Command: "java",
			Stderr:  "error details",
			Err:     baseErr,
		}
		expected := "command error [java]: exit code 1\nstderr: error details"
		if got := err.Error(); got != expected {
			t.Errorf("CommandError.Error() = %v, want %v", got, expected)
		}

		// Test Unwrap
		if got := err.Unwrap(); got != baseErr {
			t.Errorf("CommandError.Unwrap() = %v, want %v", got, baseErr)
		}
	})
}

func TestStandardErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrEmptyInput", ErrEmptyInput, "input text is empty"},
		{"ErrNilMessage", ErrNilMessage, "protobuf message cannot be nil"},
		{"ErrNoAnnotators", ErrNoAnnotators, "at least one annotator is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("%s.Error() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
