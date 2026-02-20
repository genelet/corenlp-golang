package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// HttpClient runs Stanford CoreNLP as a HTTP client
// The CoreNLP server must be actively running.
//
// see
// https://stanfordnlp.github.io/CoreNLP/index.html
type HttpClient struct {
	// a slice of annotators. e.g. []string{"tokenize","ssplit","pos","depparse"}
	Annotators []string

	// server's URL
	URL string

	// HTTPClient is the underlying HTTP client used for requests.
	// If nil, a default client with 30 second timeout will be used.
	HTTPClient *http.Client
}

// DefaultHTTPClient returns an http.Client with sensible defaults for CoreNLP.
// It includes a 30 second timeout and uses http.DefaultTransport for connection pooling.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: http.DefaultTransport,
	}
}

// NewHttpClient creates an instance of HttpClient for connecting to a CoreNLP server.
//
// Parameters:
//   - annotators: the list of annotators to run
//   - args[0], optional: the server address (default: "http://127.0.0.1:9000")
//
// Example usage with string annotators (backwards compatible):
//
//	client := NewHttpClient([]string{"tokenize", "ssplit", "pos"}, "http://localhost:9000")
//
// Example usage with type-safe Annotator constants:
//
//	client := NewHttpClientWithAnnotators([]Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS})
//
// Or use predefined combinations:
//
//	client := NewHttpClientWithAnnotators(BasicAnnotators, "http://localhost:9000")
func NewHttpClient(annotators []string, args ...string) *HttpClient {
	curl := "http://127.0.0.1:9000"
	if len(args) > 0 {
		curl = args[0]
	}

	// Fix: Check last character, not last 2 characters
	if len(curl) > 0 && curl[len(curl)-1:] != "/" {
		curl += "/"
	}
	return &HttpClient{
		Annotators: annotators,
		URL:        curl,
		HTTPClient: DefaultHTTPClient(),
	}
}

// NewHttpClientWithAnnotators creates an instance of HttpClient using type-safe Annotator constants.
// This provides better IDE autocomplete and compile-time checking for annotator names.
//
// Parameters are the same as NewHttpClient, but annotators use the Annotator type.
//
// Example:
//
//	client := NewHttpClientWithAnnotators(
//	    []Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS, AnnotatorLemma},
//	    "http://localhost:9000",
//	)
//
// Or use predefined combinations:
//
//	client := NewHttpClientWithAnnotators(BasicAnnotators)
func NewHttpClientWithAnnotators(annotators []Annotator, args ...string) *HttpClient {
	return NewHttpClient(AnnotatorsToStrings(annotators), args...)
}

// Runs on the input file, and gets the NLP data in msg
//
// Note that Document{} is the root component in the auto-generated NLP protobuf package.
func (h *HttpClient) Run(ctx context.Context, input string, msg protoreflect.ProtoMessage) error {
	data, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	return h.RunText(ctx, data, msg)
}

// RunText runs NLP analysis on the text string and populates msg with the results.
//
// The msg parameter should typically be a pointer to nlp.Document{}.
// Returns an error if the text is empty, msg is nil, or the server request fails.
func (h *HttpClient) RunText(ctx context.Context, text []byte, msg protoreflect.ProtoMessage) error {
	// Validate inputs
	if len(text) == 0 {
		return ErrEmptyInput
	}
	if msg == nil {
		return ErrNilMessage
	}

	str := ``
	if h.Annotators != nil {
		str = `"annotators":"` + strings.Join(h.Annotators, ",") + `",`
	}
	curl := h.URL + `?properties=` + url.QueryEscape(`{`+str+`"outputFormat":"serialized","serializer":"edu.stanford.nlp.pipeline.ProtobufAnnotationSerializer"}`)

	req, err := http.NewRequestWithContext(ctx, "POST", curl, bytes.NewReader(text))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	httpClient := h.HTTPClient
	if httpClient == nil {
		httpClient = DefaultHTTPClient()
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return &ServerError{
			URL:     h.URL,
			Message: err.Error(),
		}
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return &ServerError{
			URL:        h.URL,
			StatusCode: res.StatusCode,
			Message:    fmt.Sprintf("HTTP status %s", res.Status),
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return BytesUnmarshal(body, msg)
}
