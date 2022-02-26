package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Cmd runs Stanford CoreNLP process under command line.
// The original Java-based CoreNLP package must be downloaded
// and installed properly.
//
// see
// https://stanfordnlp.github.io/CoreNLP/index.html
//
type Cmd struct {
// a slice of annotators. e.g. []string{"tokenize","ssplit","pos","depparse"}
	Annotators  []string

// Classpath to the Java.
	ClassPath   string

// Class to run for annotators
	Class       string

	javaCmd     string

// extra arguments for the Java command
	Args        []string
}

// NewCmd creates an instance of Cmd
// annotators: the list of annotators.
// args[0], optional: the Java Classpath
// args[1], optional: the Java class
// args[2], optional: the Java command
// args[3:], optional: other arguments
//
// e.g. if the CoreNLP is downloaded and unzipped to /home/user/standford
// you can create an instance:
// NewCmd([]string{"tokenize","ssplit","pos"}, "/home/user/standford/*")
//
// see
// https://stanfordnlp.github.io/CoreNLP/cmdline.html
//
func NewCmd(annotators []string, args ...string) *Cmd {
	cp := "*"
	c  := "edu.stanford.nlp.pipeline.StanfordCoreNLP"
	java := "java"
	if args != nil && len(args)>0 {
		cp = args[0]
		args = args[1:]
	}
	if args != nil && len(args)>0 {
		c = args[0]
		args = args[1:]
	}
	if args != nil && len(args)>0 {
		java = args[0]
		args = args[1:]
	}

	return &Cmd{annotators, cp, c, java, args}
}

// Run on the input file, and get the parsed document in msg
//
// Document is the root component in the auto-generated protobuf GO package
// from the protobuf definition file
// https://github.com/stanfordnlp/CoreNLP/blob/main/src/edu/stanford/nlp/pipeline/CoreNLP.proto
// 
func (self *Cmd) Run(ctx context.Context, input string, msg protoreflect.ProtoMessage) error {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}
	return self.RunText(ctx, data, msg)
}

// RunText on the text string, and get the parsed document in msg
//
func (self *Cmd) RunText(ctx context.Context, text []byte, msg protoreflect.ProtoMessage) error {
	outputDir, err := ioutil.TempDir("", "coreNLP")
	if err != nil {
		return err
	}
	defer os.RemoveAll(outputDir)

	input := filepath.Join(outputDir, "input.text")
	if err = ioutil.WriteFile(input, text, 0666); err != nil {
		return err
	}

	args := self.Args
	if self.ClassPath != "" {
		args = append(args, "-cp", self.ClassPath)
	}
	args = append(args, self.Class)
	if self.Annotators != nil && len(self.Annotators) > 0 {
		args = append(args, "-annotators", strings.Join(self.Annotators, ","))
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

	cmd := exec.CommandContext(ctx, self.javaCmd, args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %s", err.Error(), stderr.String())
	}

	data, err := ioutil.ReadFile(input+".ser.gz")
	if err != nil {
		return err
	}

	return BytesUnmarshal(data, msg)
}
