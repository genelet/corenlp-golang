# corenlp-golang

[![GoDoc](https://godoc.org/github.com/genelet/corenlp-golang?status.svg)](https://godoc.org/github.com/genelet/corenlp-golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/genelet/corenlp-golang)](https://goreportcard.com/report/github.com/genelet/corenlp-golang)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.17+-lightblue.svg)](https://go.dev/)

A powerful, type-safe Go client for [Stanford CoreNLP](https://stanfordnlp.github.io/CoreNLP/index.html), providing access to the complete NLP data set defined by CoreNLP.proto.

## ✨ Features

- 🔒 **Type-Safe Annotators** - Compile-time checking with IDE autocomplete
- 🎯 **Predefined Pipelines** - Ready-to-use annotator combinations  
- 🛠️ **Helper Functions** - Easy extraction of tokens, entities, lemmas, and more
- ⚡ **Two Client Modes** - Command-line or HTTP server
- 🔄 **Interface-Based** - Easy testing and dependency injection
- 🎨 **Custom Error Types** - Better error handling with `errors.As()`
- 📦 **Zero Dependencies** - Only requires `google.golang.org/protobuf`
- ✅ **Fully Tested** - Comprehensive test coverage

## 📋 Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Command Line Client](#command-line-client)
  - [HTTP Client](#http-client)
  - [Type-Safe Annotators](#type-safe-annotators-recommended)
- [Helper Functions](#helper-functions)
- [Annotator Reference](#annotator-reference)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## 🚀 Installation

```bash
go get -u github.com/genelet/corenlp-golang
```

**Current Version**: v0.5.10 (compatible with Stanford CoreNLP 4.5.10)

### Prerequisites

Choose one of the following:

**Option 1: Command Line** - Download and install [Stanford CoreNLP](https://stanfordnlp.github.io/CoreNLP/download.html)

**Option 2: HTTP Service** - Run CoreNLP as a server (recommended for production)

## ⚡ Quick Start

### HTTP Client (Recommended)

```go
package main

import (
    "context"
    "fmt"
    "github.com/genelet/corenlp-golang/client"
    "github.com/genelet/corenlp-golang/nlp"
)

func main() {
    // Create client with predefined annotators
    c := client.NewHttpClientWithAnnotators(
        client.BasicAnnotators,
        "http://localhost:9000",
    )
    
    // Analyze text
    doc := &nlp.Document{}
    text := "Stanford University is located in California."
    err := c.RunText(context.Background(), []byte(text), doc)
    if err != nil {
        panic(err)
    }
    
    // Extract entities
    entities := client.ExtractNamedEntities(doc)
    fmt.Printf("Entities: %v\n", entities)
    // Output: Entities: map[LOCATION:[California] ORGANIZATION:[Stanford University]]
}
```

### Command Line Client

```go
func main() {
    // Create command-line client
    cmd := client.NewCmdWithAnnotators(
        client.BasicAnnotators,
        "/home/user/stanford-corenlp-4.5.10/*",
    )
    
    doc := &nlp.Document{}
    err := cmd.RunText(context.Background(), []byte("Your text here."), doc)
    if err != nil {
        panic(err)
    }
}
```

## 📖 Usage

### Comparison: Cmd vs HttpClient

| Feature | `Cmd` (Command Line) | `HttpClient` (HTTP) |
|---------|---------------------|---------------------|
| **Setup** | Download CoreNLP zip | Run CoreNLP server |
| **Best For** | Development, testing | Production, high-volume |
| **Resource Usage** | Spawns Java process per request | Reuses server connection |
| **Performance** | Slower (process overhead) | Faster (persistent server) |
| **Scalability** | Limited | High (server handles pooling) |
| **Dependencies** | Local CoreNLP installation | Running CoreNLP server |

### Command Line Client

#### Setup

1. Download [Stanford CoreNLP](https://stanfordnlp.github.io/CoreNLP/download.html)
2. Unzip to a directory (e.g., `/home/user/stanford-corenlp-4.5.10/`)
3. Verify installation:

```bash
cd /home/user/stanford-corenlp-4.5.10/
java -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLP -file input.txt
```

#### Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/genelet/corenlp-golang/client"
    "github.com/genelet/corenlp-golang/nlp"
)

func main() {
    // Create client
    cmd := client.NewCmd(
        []string{"tokenize", "ssplit", "pos", "lemma"},
        "/home/user/stanford-corenlp-4.5.10/*",
    )
    
    // Analyze text
    doc := &nlp.Document{}
    text := `Stanford University is located in California. It is a great university, founded in 1891.`
    err := cmd.RunText(context.Background(), []byte(text), doc)
    if err != nil {
        panic(err)
    }
    
    // Print results
    fmt.Printf("%12s %12s %8s\n", "Word", "Lemma", "POS")
    fmt.Println("  --------------------------------")
    for _, token := range doc.Sentence[0].Token {
        fmt.Printf("%12s %12s %8s\n", *token.Word, *token.Lemma, *token.Pos)
    }
}
```

**Output:**
```
        Word        Lemma      POS
  --------------------------------
     Stanford     Stanford      NNP
   University   University      NNP
           is           be      VBZ
      located       locate      VBN
           in           in       IN
   California   California      NNP
            .            .        .
```

### HTTP Client

#### Setup

Run CoreNLP as a server (recommended for production):

**Using Docker:**
```bash
docker run -p 9000:9000 --name corenlp --rm \
  nlpbox/corenlp:latest
```

**Using systemd (Linux):**
```bash
# Create service file
cat > ~/.config/systemd/user/coreNLP.service << EOF
[Unit]
Description=CoreNLP Server at 9000

[Service]
Type=simple
WorkingDirectory=/home/user/stanford-corenlp-4.5.10
Environment=CLASSPATH=/home/user/stanford-corenlp-4.5.10/*:
ExecStart=/usr/bin/java -mx4g -cp "/home/user/stanford-corenlp-4.5.10/*" edu.stanford.nlp.pipeline.StanfordCoreNLPServer -port 9000 -timeout 15000

[Install]
WantedBy=default.target
EOF

# Enable and start
systemctl --user enable coreNLP.service
systemctl --user start coreNLP.service
```

**Manual start:**
```bash
cd /home/user/stanford-corenlp-4.5.10/
java -mx4g -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLPServer -port 9000 -timeout 15000
```

#### Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/genelet/corenlp-golang/client"
    "github.com/genelet/corenlp-golang/nlp"
)

func main() {
    // Create HTTP client
    c := client.NewHttpClient(
        []string{"tokenize", "ssplit", "pos", "lemma", "ner"},
        "http://localhost:9000",
    )
    
    doc := &nlp.Document{}
    text := `Stanford University is located in California.`
    err := c.RunText(context.Background(), []byte(text), doc)
    if err != nil {
        panic(err)
    }
    
    // Use helper to extract entities
    entities := client.ExtractNamedEntities(doc)
    for entityType, names := range entities {
        fmt.Printf("%s: %v\n", entityType, names)
    }
}
```

**Output:**
```
ORGANIZATION: [Stanford University]
LOCATION: [California]
```

### Type-Safe Annotators (Recommended)

Use type-safe constants for better IDE support and compile-time checking:

```go
// ✅ Type-safe with autocomplete
c := client.NewHttpClientWithAnnotators(
    []client.Annotator{
        client.AnnotatorTokenize,
        client.AnnotatorSSplit,
        client.AnnotatorPOS,
        client.AnnotatorLemma,
        client.AnnotatorNER,
    },
    "http://localhost:9000",
)

// ❌ Old way (still supported but not recommended)
c := client.NewHttpClient(
    []string{"tokenize", "ssplit", "pos", "lemma", "ner"},
    "http://localhost:9000",
)
```

#### Predefined Annotator Combinations

Save time with ready-to-use pipelines:

```go
// Basic text processing
client.BasicAnnotators
// → tokenize, ssplit, pos, lemma

// Syntax analysis
client.SyntaxAnnotators  
// → BasicAnnotators + parse, depparse

// Named Entity Recognition
client.NERAnnotators
// → BasicAnnotators + ner, entitymentions

// Semantic analysis (full pipeline)
client.SemanticAnnotators
// → tokenize, ssplit, pos, lemma, ner, parse, depparse, coref

// Relation extraction
client.RelationExtractionAnnotators
// → tokenize, ssplit, pos, lemma, depparse, natlog, openie
```

**Example:**
```go
// Quick setup for NER tasks
c := client.NewHttpClientWithAnnotators(
    client.NERAnnotators,
    "http://localhost:9000",
)
```

## 🛠️ Helper Functions

Convenient utilities to extract common NLP data:

```go
import "github.com/genelet/corenlp-golang/client"

// Assume doc is populated from running NLP analysis
doc := &nlp.Document{}
c.RunText(ctx, text, doc)

// Extract tokens as strings
tokens := client.ExtractTokens(doc)
// ["Stanford", "University", "is", "located", "in", "California", "."]

// Extract sentences
sentences := client.ExtractSentences(doc)
// ["Stanford University is located in California ."]

// Extract named entities (grouped by type)
entities := client.ExtractNamedEntities(doc)
// {"ORGANIZATION": ["Stanford University"], "LOCATION": ["California"]}

// Extract lemmas (base forms)
lemmas := client.ExtractLemmas(doc)
// ["Stanford", "University", "be", "locate", "in", "California", "."]

// Extract part-of-speech tags
posTags := client.ExtractPOSTags(doc)
// ["NNP", "NNP", "VBZ", "VBN", "IN", "NNP", "."]

// Extract tokens with all metadata
tokensWithMeta := client.ExtractTokensWithMetadata(doc)
// []TokenWithMetadata{{Word:"Stanford", Lemma:"Stanford", POS:"NNP", NER:"ORGANIZATION"}, ...}
```

## 📚 Annotator Reference

### Core Annotators

| Annotator | Type-Safe Constant | Description | Dependencies |
|-----------|-------------------|-------------|--------------|
| **tokenize** | `AnnotatorTokenize` | Splits text into tokens/words | - |
| **cleanxml** | `AnnotatorCleanXML` | Removes XML tags from document | - |
| **ssplit** | `AnnotatorSSplit` | Splits text into sentences | tokenize |
| **docdate** | `AnnotatorDocDate` | Extracts document date | - |
| **pos** | `AnnotatorPOS` | Part-of-speech tagging | tokenize, ssplit |
| **lemma** | `AnnotatorLemma` | Lemmatization (base forms) | tokenize, ssplit, pos |

### Named Entity Recognition

| Annotator | Type-Safe Constant | Description | Dependencies |
|-----------|-------------------|-------------|--------------|
| **ner** | `AnnotatorNER` | Named entity recognition | tokenize, ssplit, pos, lemma |
| **regexner** | `AnnotatorRegexNER` | Rule-based NER with regex | tokenize, ssplit, pos, lemma, ner |
| **entitymentions** | `AnnotatorEntityMentions` | Entity mention detection | tokenize, ssplit, pos, lemma, ner |
| **entitylink** | `AnnotatorEntityLink` | Link entities to Wikipedia | tokenize, ssplit, pos, lemma, ner, entitymentions |

### Parsing & Syntax

| Annotator | Type-Safe Constant | Description | Dependencies |
|-----------|-------------------|-------------|--------------|
| **parse** | `AnnotatorParse` | Constituency parsing | tokenize, ssplit, pos |
| **depparse** | `AnnotatorDepparse` | Dependency parsing | tokenize, ssplit, pos |

### Coreference & Semantics

| Annotator | Type-Safe Constant | Description | Dependencies |
|-----------|-------------------|-------------|--------------|
| **coref** | `AnnotatorCoref` | Coreference resolution | tokenize, ssplit, pos, lemma, ner, parse |
| **sentiment** | `AnnotatorSentiment` | Sentiment analysis | tokenize, ssplit, parse |
| **natlog** | `AnnotatorNatlog` | Natural logic semantics | tokenize, ssplit, pos, lemma, depparse |
| **openie** | `AnnotatorOpenie` | Open information extraction | tokenize, ssplit, pos, lemma, depparse, natlog |
| **truecase** | `AnnotatorTruecase` | Determines true case of tokens | - |
| **udfeats** | `AnnotatorUDFeats` | Universal Dependencies features | tokenize, ssplit, pos |

### Information Extraction

| Annotator | Type-Safe Constant | Description | Dependencies |
|-----------|-------------------|-------------|--------------|
| **relation** | `AnnotatorRelation` | Relation extraction | tokenize, ssplit, pos, lemma, ner, parse |
| **kbp** | `AnnotatorKBP` | Knowledge base population | tokenize, ssplit, pos, lemma, ner, parse, coref |

### Text Processing

| Annotator | Type-Safe Constant | Description | Dependencies |
|-----------|-------------------|-------------|--------------|
| **quote** | `AnnotatorQuote` | Quote extraction | tokenize, ssplit |
| **quote.attribution** | `AnnotatorQuoteAttribution` | Quote attribution to speakers | tokenize, ssplit, pos, lemma, ner, depparse, coref, quote |
| **tokensregex** | `AnnotatorTokensRegex` | TokensRegex pattern matching | tokenize, ssplit |

See [annotators.go](client/annotators.go) for the complete list of all 25+ available annotators with detailed documentation.

## 💡 Best Practices

### Error Handling

Use custom error types for better error handling:

```go
import (
    "errors"
    "github.com/genelet/corenlp-golang/client"
)

err := c.RunText(ctx, text, doc)
if err != nil {
    // Check for server errors
    var serverErr *client.ServerError
    if errors.As(err, &serverErr) {
        log.Printf("Server %s failed with status %d: %s", 
            serverErr.URL, serverErr.StatusCode, serverErr.Message)
        return
    }
    
    // Check for command execution errors
    var cmdErr *client.CommandError
    if errors.As(err, &cmdErr) {
        log.Printf("Command %s failed: %v\nStderr: %s", 
            cmdErr.Command, cmdErr.Err, cmdErr.Stderr)
        return
    }
    
    // Check for standard errors
    if errors.Is(err, client.ErrEmptyInput) {
        log.Println("Input text is empty")
        return
    }
    
    // Generic error
    log.Printf("NLP processing failed: %v", err)
}
```

### Interface-Based Design

Both `Cmd` and `HttpClient` implement the `Client` interface for flexibility:

```go
type NLPProcessor struct {
    client client.Client
}

func NewNLPProcessor(useHTTP bool, config string) *NLPProcessor {
    var c client.Client
    
    if useHTTP {
        c = client.NewHttpClientWithAnnotators(
            client.NERAnnotators,
            config, // server URL
        )
    } else {
        c = client.NewCmdWithAnnotators(
            client.NERAnnotators,
            config, // classpath
        )
    }
    
    return &NLPProcessor{client: c}
}

func (p *NLPProcessor) Process(text string) (*nlp.Document, error) {
    doc := &nlp.Document{}
    err := p.client.RunText(context.Background(), []byte(text), doc)
    return doc, err
}
```

### Performance Tips

1. **Use HTTP client in production** - Better resource management and connection pooling
2. **Reuse client instances** - Don't create new clients for each request
3. **Choose minimal annotators** - Only include what you need; fewer annotators = faster processing
4. **Use predefined combinations** - They're optimized for common use cases
5. **Batch processing** - For large datasets, run your own CoreNLP server
6. **Context with timeout** - Prevent hanging requests:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := client.RunText(ctx, text, doc)
```

## 🔧 Troubleshooting

### Common Issues

#### "Could not find or load main class edu.stanford.nlp.pipeline.StanfordCoreNLP"

**Solution:** Check your classpath. Make sure it points to the correct CoreNLP directory:

```go
// ❌ Wrong
cmd := client.NewCmd(annotators, "/home/user/stanford-corenlp-4.5.10")

// ✅ Correct - include the wildcard
cmd := client.NewCmd(annotators, "/home/user/stanford-corenlp-4.5.10/*")
```

#### "HTTP status 404" or "Connection refused"

**Solution:** Ensure the CoreNLP server is running:

```bash
# Check if server is running
curl http://localhost:9000

# Start server if not running
java -mx4g -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLPServer -port 9000
```

#### "OutOfMemoryError" from Java

**Solution:** Increase Java heap size:

```bash
# Increase from 4GB to 8GB
java -mx8g -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLPServer -port 9000
```

#### Slow performance

**Solutions:**
- Use fewer annotators (only what you need)
- Switch from `Cmd` to `HttpClient`
- Increase server memory with `-mx` flag
- Run server on dedicated hardware
- Process in batches

#### Empty or nil results

**Solution:** Check that required annotators are included. For example, to extract entities:

```go
// ❌ Missing NER annotator
client.NewHttpClient([]string{"tokenize", "ssplit"}, url)

// ✅ Include NER
client.NewHttpClientWithAnnotators(client.NERAnnotators, url)
```

## 📝 Protocol Buffer Definition

The NLP data structure is defined by [CoreNLP.proto](https://github.com/stanfordnlp/CoreNLP/blob/main/src/edu/stanford/nlp/pipeline/CoreNLP.proto).

The auto-generated Go package is at [github.com/genelet/corenlp-golang/nlp](https://github.com/genelet/corenlp-golang/tree/main/nlp).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development

```bash
# Clone repository
git clone https://github.com/genelet/corenlp-golang.git
cd corenlp-golang

# Run tests
go test ./client/... -v

# Run linter
go vet ./...

# Format code
go fmt ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Documentation**: [https://godoc.org/github.com/genelet/corenlp-golang](https://godoc.org/github.com/genelet/corenlp-golang)
- **Stanford CoreNLP**: [https://stanfordnlp.github.io/CoreNLP/](https://stanfordnlp.github.io/CoreNLP/)
- **Issues**: [https://github.com/genelet/corenlp-golang/issues](https://github.com/genelet/corenlp-golang/issues)

## ⭐ Star History

If you find this project helpful, please consider giving it a star!

---

**Made with ❤️ by the CoreNLP-Golang community**
