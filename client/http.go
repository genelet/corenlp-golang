package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// HttpClient runs Stanford CoreNLP process as a HTTP client
// The CoreNLP server must be actively running.
//
// see
// https://stanfordnlp.github.io/CoreNLP/index.html
//
type HttpClient struct {
// a slice of annotators. e.g. []string{"tokenize","ssplit","pos","depparse"}
	Annotators []string

// server's URL 
	URL        string
}

// NewHttpClient creates an instance of HttpClient
// annotators: the list of annotators.
// args[0], optional: the server address, default to http://127.0.0.1:9000
//
//
func NewHttpClient(annotators []string, args ...string) *HttpClient {
    curl := "http://127.0.0.1:9000"
	if args != nil {
		curl = args[0]
	}

	if curl[len(curl)-2:] != `/` {
		curl += `/`
	}
	return &HttpClient{annotators, curl}
}

// Run using the input file, and get the parsed document in msg
//
// Document is the root component in the auto-generated protobuf GO package
// from the protobuf definition file
// https://github.com/stanfordnlp/CoreNLP/blob/main/src/edu/stanford/nlp/pipeline/CoreNLP.proto
// 
func (self *HttpClient) Run(ctx context.Context, input string, msg protoreflect.ProtoMessage) error {
    data, err := ioutil.ReadFile(input)
    if err != nil {
        return err
    }
    return self.RunText(ctx, data, msg)
}

// RunText on the text string, and get the parsed document in msg
//
func (self *HttpClient) RunText(ctx context.Context, text []byte, msg protoreflect.ProtoMessage) error {
	str := ``
	if self.Annotators != nil {
		str = `"annotators":"` + strings.Join(self.Annotators, ",") + `",`
	}
	curl := self.URL + `?properties=`+ url.QueryEscape(`{`+str+`"outputFormat":"serialized","serializer":"edu.stanford.nlp.pipeline.ProtobufAnnotationSerializer"}`)

	req, err := http.NewRequestWithContext(ctx, "POST", curl, bytes.NewReader(text))
	if err != nil {
		return err
	}

	defaultClient := &http.Client{Transport: http.DefaultTransport}
	res, err := defaultClient.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("HTTP status %s\n", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	return BytesUnmarshal(body, msg)
}
