package client

import (
	"context"
	"testing"

	"github.com/genelet/coreNLP/nlp"
)

func TestCmd(t *testing.T) {
	cmd := NewCmd([]string{"tokenize","ssplit","pos","lemma","parse","depparse"}, "/home/peter/stanford-corenlp-4.4.0/*")

	pb := &nlp.Document{}
	var text = `President Xi Jinping of Chaina, on his first state visit to the United States, showed off his familiarity with American history and pop culture on Tuesday night.`

	err := cmd.RunText(context.Background(), []byte(text), pb)
	if err != nil { t.Fatal(err) }

	if pb.String()[:168] != `text:"President Xi Jinping of Chaina, on his first state visit to the United States, showed off his familiarity with American history and pop culture on Tuesday night."` {
		t.Errorf("%s", pb.String()[:168])
    }
}
