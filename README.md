# corenlp-golang
*corenlp-golang* is a GO client to access the complete set of [Stanford CoreNLP](https://stanfordnlp.github.io/CoreNLP/index.html) components.

[![GoDoc](https://godoc.org/github.com/genelet/corenlp-golang?status.svg)](https://godoc.org/github.com/genelet/corenlp-golang)

## 1. Installation

> $ go get -u github.com/genelet/corenlp-golang

This GO client package should be used with either command line program or web service.

#### 1.1) Command Line

Download the Stanford CoreNLP and unzip it:

[https://stanfordnlp.github.io/CoreNLP/download.html](https://stanfordnlp.github.io/CoreNLP/download.html).

Go to the directory and make sure to the following command line can run properly:

```bash
$ java -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLP -file input.txt
```

#### 1.2) Web Service

CoreNLP can also be launched as a http Web service. 

For example, under a Ubuntu account, create the following startup script and run it as a service.

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

The data that CoreNLP returns from neutral language processing can be summaried in a [protocol buffer](https://developers.google.com/protocol-buffers/docs/overview):

https://github.com/stanfordnlp/CoreNLP/blob/main/src/edu/stanford/nlp/pipeline/CoreNLP.proto

The auto-generated GO packge is included in [github.com/genelet/corenlp-golang/nlp](https://github.com/genelet/corenlp-golang/tree/main/nlp)

<br /><br />

## 2. Usage

There are two functions implemented:

- *Run*: it reads neutral language from a text file, and returns the NLP data as protobuf message.
- *RunText*: the same as *Run* but reads text directly.

#### 2.1) Use Command Line Interface

```go
package main

import (
    "context"
    "fmt"
    "github.com/genelet/corenlp-golang/client"
    "github.com/genelet/corenlp-golang/nlp"
)

func main() {
    // assuming the Stanford CoreNLP is downloaded into /home/user/stanford-corenlp-4.4.0
    // create a new Cmd instance
    cmd := client.NewCmd([]string{"tokenize","ssplit","pos","lemma","parse","depparse"}, "/home/user/stanford-corenlp-4.4.0/*")

    // a reference to the nlp Document
    pb := &nlp.Document{}

    // run NLP and receive data in pb
    err := cmd.RunText(context.Background(), []byte(`Stanford University is located in California. It is a great university, founded in 1891.`), pb)
    if err != nil { panic(err) }

    // print some result
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

#### 2.2) Web Service

Using the web service is almost identical:

```go
package main

import (
    "context"
    "fmt"
    "github.com/genelet/corenlp-golang/client"
    "github.com/genelet/corenlp-golang/nlp"
)

func main() {
    // assuming the Stanford CoreNLP is running at http://localhost:9000
    // create a new HttpClient instance
    cmd := client.NewHttpClient([]string{"tokenize","ssplit","pos","lemma","parse","depparse"}, "http://localhost:9000")

    // a reference to the nlp Document
    pb := &nlp.Document{}

    // run NLP and receive data in pb
    err := cmd.RunText(context.Background(), []byte(`Stanford University is located in California. It is a great university, founded in 1891.`), pb)
    if err != nil { panic(err) }
    
    // print some result
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

Please check [https://godoc.org/github.com/genelet/corenlp-golang](https://godoc.org/github.com/genelet/corenlp-golang) for the complete document.
