package client

import (
	"context"
	"testing"

	"github.com/genelet/corenlp-golang/nlp"
)

func TestCoreNLP(t *testing.T) {
	client := NewHttpClient([]string{"tokenize", "ssplit", "pos", "lemma", "parse", "depparse"})
	pb := &nlp.Document{}
	var text = `President Xi Jinping of Chaina, on his first state visit to the United States, showed off his familiarity with American history and pop culture on Tuesday night.`
	err := client.RunText(context.Background(), []byte(text), pb)
	if err != nil {
		t.Fatal(err)
	}

	if pb.String()[:168] != `text:"President Xi Jinping of Chaina, on his first state visit to the United States, showed off his familiarity with American history and pop culture on Tuesday night."` {
		t.Errorf("%s", pb.String()[:168])
	}
}
