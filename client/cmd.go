package client

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Cmd runs Stanford CoreNLP under command line.
// The original Java-based CoreNLP package must be downloaded
// and installed properly.
//
// see
// https://stanfordnlp.github.io/CoreNLP/index.html
type Cmd struct {
	// a slice of annotators. e.g. []string{"tokenize","ssplit","pos","depparse"}
	Annotators []string

	// Classpath to the Java.
	ClassPath string

	// Class to run for annotators
	Class string

	javaCmd string

	// extra arguments for the Java command
	Args []string
}

// NewCmd creates an instance of Cmd for running CoreNLP via command line.
//
// Parameters:
//   - annotators: the list of annotators to run (e.g., tokenize, ssplit, pos, lemma, ner, parse)
//   - args[0], optional: the Java Classpath (default: "*")
//   - args[1], optional: the Java class to run (default: "edu.stanford.nlp.pipeline.StanfordCoreNLP")
//   - args[2], optional: the Java command (default: "java")
//   - args[3:], optional: additional Java arguments
//
// Example usage with string annotators (backwards compatible):
//
//	cmd := NewCmd([]string{"tokenize", "ssplit", "pos"}, "/home/user/stanford-corenlp-4.5.4/*")
//
// Example usage with type-safe Annotator constants:
//
//	cmd := NewCmdWithAnnotators([]Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS}, "/home/user/stanford-corenlp-4.5.4/*")
//
// Or use predefined combinations:
//
//	cmd := NewCmdWithAnnotators(BasicAnnotators, "/home/user/stanford-corenlp-4.5.4/*")
//
// See https://stanfordnlp.github.io/CoreNLP/cmdline.html for more information.
func NewCmd(annotators []string, args ...string) *Cmd {
	cp := "*"
	c := "edu.stanford.nlp.pipeline.StanfordCoreNLP"
	java := "java"
	if len(args) > 0 {
		cp = args[0]
		args = args[1:]
	}
	if len(args) > 0 {
		c = args[0]
		args = args[1:]
	}
	if len(args) > 0 {
		java = args[0]
		args = args[1:]
	}

	return &Cmd{annotators, cp, c, java, args}
}

// NewCmdWithAnnotators creates an instance of Cmd using type-safe Annotator constants.
// This provides better IDE autocomplete and compile-time checking for annotator names.
//
// Parameters are the same as NewCmd, but annotators use the Annotator type.
//
// Example:
//
//	cmd := NewCmdWithAnnotators(
//	    []Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS, AnnotatorLemma},
//	    "/home/user/stanford-corenlp-4.5.4/*",
//	)
//
// Or use predefined combinations:
//
//	cmd := NewCmdWithAnnotators(BasicAnnotators, "/home/user/stanford-corenlp-4.5.4/*")
func NewCmdWithAnnotators(annotators []Annotator, args ...string) *Cmd {
	return NewCmd(AnnotatorsToStrings(annotators), args...)
}

// Runs on the input file, and gets the NLP data in msg.
//
// Note that Document{} is the root component in the auto-generated NLP protobuf package.
func (c *Cmd) Run(ctx context.Context, input string, msg protoreflect.ProtoMessage) error {
	data, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	return c.RunText(ctx, data, msg)
}

// RunText runs NLP analysis on the text string and populates msg with the results.
//
// The msg parameter should typically be a pointer to nlp.Document{}.
// Returns an error if the text is empty, msg is nil, or the command execution fails.
func (c *Cmd) RunText(ctx context.Context, text []byte, msg protoreflect.ProtoMessage) error {
	// Validate inputs
	if len(text) == 0 {
		return ErrEmptyInput
	}
	if msg == nil {
		return ErrNilMessage
	}

	outputDir, err := os.MkdirTemp("", "coreNLP")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(outputDir)

	input := filepath.Join(outputDir, "input.text")
	if err = os.WriteFile(input, text, 0666); err != nil {
		return fmt.Errorf("failed to write input file: %w", err)
	}

	args := c.Args
	if c.ClassPath != "" {
		args = append(args, "-cp", c.ClassPath)
	}
	args = append(args, c.Class)
	if len(c.Annotators) > 0 {
		args = append(args, "-annotators", strings.Join(c.Annotators, ","))
	}

	args = append(args,
		"-file",
		input,
		"--outputDirectory",
		outputDir,
		"-outputFormat",
		"serialized",
		"-outputSerializer",
		"edu.stanford.nlp.pipeline.ProtobufAnnotationSerializer")

	cmd := exec.CommandContext(ctx, c.javaCmd, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return &CommandError{
			Command: c.javaCmd,
			Stderr:  stderr.String(),
			Err:     err,
		}
	}

	data, err := os.ReadFile(input + ".ser.gz")
	if err != nil {
		return fmt.Errorf("failed to read output file: %w", err)
	}

	return BytesUnmarshal(data, msg)
}
