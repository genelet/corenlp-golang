# coreNLP
*coreNLP* is a Golang wrapper to access the full Standfor CoreNLP components

[![GoDoc](https://godoc.org/github.com/genelet/coreNLP?status.svg)](https://godoc.org/github.com/genelet/coreNLP)

## Installation

> $ go get -u github.com/genelet/coreNLP

### Command Line

The Stanford CoreNLP should be downloaded and install properly from

[https://stanfordnlp.github.io/CoreNLP/download.html](https://stanfordnlp.github.io/CoreNLP/download.html).

### Web Service

Instead of the command line, you may launch CoreNLP as a Web service. For example, in your personal Ubuntu account, 

1) create ~/.config/systemd/user/coreNLP.service with the following content

[Unit]
Description=CoreNLP Server at 9000

[Service]
Type=simple
WorkingDirectory=/home/user/stanford-corenlp-4.4.0
Environment=CLASSPATH=/home/user/stanford-corenlp-4.4.0/*:
ExecStart=/usr/bin/java -mx4g -cp "/home/user/stanford-corenlp-4.4.0/*" edu.stanford.nlp.pipeline.StanfordCoreNLPServer -port 9000 -timeout 15000

[Install]
WantedBy=default.target

2) systemctl --user enable coreNLP.service

3) systemctl --user start coreNLP.service

4) systemctl --user daemon-reload

### protobuf

The annotating components can be summaried in [protobuf](https://developers.google.com/protocol-buffers/docs/overview), the definition is at

https://github.com/stanfordnlp/CoreNLP/blob/main/src/edu/stanford/nlp/pipeline/CoreNLP.proto

The auto-generated packge is named as [github.com/genelet/coreNLP/nlp](github.com/genelet/coreNLP/nlp).

## Usage

Function $Run$ is implemented in both 
the command line interface and the http web interface. By assign an input text
in a file, this function returns the annotation data as a protobuf message.

### Command Line Interface

package main

import (
	"context"

	"github.com/genelet/coreNLP/client"
	"github.com/genelet/coreNLP/nlp"
)

func main() {
	// assuming the Stanford CoreNLP is downloaded into /home/user/stanford-corenlp-4.4.0

	// create a new Cmd instance
	cmd := NewCmd([]string{"tokenize","ssplit","pos","lemma","parse","depparse"}, "/home/user/stanford-corenlp-4.4.0/*")

	// a reference to the nlp Document
	pb := &nlp.Document{}

	err := cmd.RunText(context.Background(), []byte(`Stanford University is located in California. It is a great university, founded in 1891.`), pb)
    if err != nil { panic(err) }

    fmt.Printf("%12.12s %12.12s %8.8s\n", "Word", "Lemma", "Pos")
    fmt.Printf("%s\n", "  --------------------------------")
    for _, token := range pb.Sentence[0].Token {
        fmt.Printf("%12.12s %12.12s %8.8s\n", *token.Word, *token.Lemma, *token.Pos)
    }
}


        Word        Lemma      Pos
  --------------------------------
    Stanford     Stanford      NNP
  University   University      NNP
          is           be      VBZ
     located       locate      VBN
          in           in       IN
  California   California      NNP
           .            .        .

### Web Service

package main

import (
    "context"

    "github.com/genelet/coreNLP/client"
    "github.com/genelet/coreNLP/nlp"
)

func main() {
    // assuming the Stanford CoreNLP is downloaded into /home/user/stanford-corenlp-4.4.0

    // create a new Cmd instance
    cmd := NewCmd([]string{"tokenize","ssplit","pos","lemma","parse","depparse"}, "/home/user/stanford-corenlp-4.4.0/*")

    // a reference to the nlp Document
    pb := &nlp.Document{}

    err := cmd.RunText(context.Background(), []byte(`Stanford University is located in California. It is a great university, founded in 1891.`), pb)
    if err != nil { panic(err) }

    fmt.Printf("%12.12s %12.12s %8.8s\n", "Word", "Lemma", "Pos")
    fmt.Printf("%s\n", "  --------------------------------")
    for _, token := range pb.Sentence[0].Token {
        fmt.Printf("%12.12s %12.12s %8.8s\n", *token.Word, *token.Lemma, *token.Pos)
    }
}


        Word        Lemma      Pos
  --------------------------------
    Stanford     Stanford      NNP
  University   University      NNP
          is           be      VBZ
     located       locate      VBN
          in           in       IN
  California   California      NNP
           .            .        .

Please check [here]() for the complete docuemnt
