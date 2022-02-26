# coreNLP
*coreNLP* is a Golang wrapper to access the full Stanford CoreNLP components.

[![GoDoc](https://godoc.org/github.com/genelet/coreNLP?status.svg)](https://godoc.org/github.com/genelet/coreNLP)

## 1. Installation

> $ go get -u github.com/genelet/coreNLP

This GO client package should be used with either the CoreNLP command line program or a CoreNLP web service. Please check the following document for detail:

[https://stanfordnlp.github.io/CoreNLP/index.html](https://stanfordnlp.github.io/CoreNLP/index.html)


#### 1.1) Command Line

Download the Stanford CoreNLP here and install it properly:

[https://stanfordnlp.github.io/CoreNLP/download.html](https://stanfordnlp.github.io/CoreNLP/download.html).


#### 1.2) Web Service

Besides the command line, you may launch CoreNLP as a Web service. 

For example, in your personal Ubuntu account, 
```bash
$ vi ~/.config/systemd/user/coreNLP.service 

[Unit]
Description=CoreNLP Server at 9000

[Service]
Type=simple
WorkingDirectory=/home/user/stanford-corenlp-4.4.0
Environment=CLASSPATH=/home/user/stanford-corenlp-4.4.0/*:
ExecStart=/usr/bin/java -mx4g -cp "/home/user/stanford-corenlp-4.4.0/*" edu.stanford.nlp.pipeline.StanfordCoreNLPServer -port 9000 -timeout 15000

[Install]
WantedBy=default.target

$ systemctl --user enable coreNLP.service
$ systemctl --user start coreNLP.service
$ systemctl --user daemon-reload
```

#### 1.3) The *proto* definition

The annotating components can be summaried in [protobuf](https://developers.google.com/protocol-buffers/docs/overview):

https://github.com/stanfordnlp/CoreNLP/blob/main/src/edu/stanford/nlp/pipeline/CoreNLP.proto

The auto-generated GO packge is named as [github.com/genelet/coreNLP/nlp](github.com/genelet/coreNLP/nlp).

<br /><br />

## 2. Usage

Function *Run* is implemented in both 
the command line interface and the http web interface. By assign an input text
in a file, this function returns the NLP data as a protobuf message.

Function *RunText* is the same program but using text input directly.

#### 2.1) Command Line Interface

```go
package main

import (
	"context"
        "fmt"
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
```
It outputs:
```bash
        Word        Lemma      Pos
  --------------------------------
    Stanford     Stanford      NNP
  University   University      NNP
          is           be      VBZ
     located       locate      VBN
          in           in       IN
  California   California      NNP
           .            .        .
```

#### 2.2) Web Service

Using the web service is almost identical:

```go
package main

import (
    "context"
    "fmt"
    "github.com/genelet/coreNLP/client"
    "github.com/genelet/coreNLP/nlp"
)

func main() {
    // assuming the Stanford CoreNLP is running at http://localhost:9000

    // create a new HttpClient instance
    cmd := NewHttpClient([]string{"tokenize","ssplit","pos","lemma","parse","depparse"}, "http://localhost:9000")

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
```

It outputs:
```
        Word        Lemma      Pos
  --------------------------------
    Stanford     Stanford      NNP
  University   University      NNP
          is           be      VBZ
     located       locate      VBN
          in           in       IN
  California   California      NNP
           .            .        .
```

Please check [here]() for the complete document.
