package client

import (
	"context"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Client is the common interface implemented by both Cmd and HttpClient.
// This interface allows for easier testing, mocking, and dependency injection.
//
// Example usage:
//
//	var client Client
//	if useHTTP {
//	    client = NewHttpClient(annotators, serverURL)
//	} else {
//	    client = NewCmd(annotators, classPath)
//	}
//	err := client.RunText(ctx, text, &nlp.Document{})
type Client interface {
	// Run processes text from a file and populates the protobuf message with NLP results.
	// The msg parameter should typically be a pointer to nlp.Document{}.
	Run(ctx context.Context, input string, msg protoreflect.ProtoMessage) error

	// RunText processes text directly and populates the protobuf message with NLP results.
	// The msg parameter should typically be a pointer to nlp.Document{}.
	RunText(ctx context.Context, text []byte, msg protoreflect.ProtoMessage) error
}
