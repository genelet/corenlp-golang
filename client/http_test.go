package client

import (
	"context"
	"testing"

	"github.com/genelet/corenlp-golang/nlp"
)

func TestCoreNLPReturnsExpectedTokenAmount(t *testing.T) {
	tests := map[string]struct {
		given string
		want  int
	}{
		"state visit sentence": {
			given: "President Xi Jinping of Chaina, on his first state visit to the United States, showed off his familiarity with American history and pop culture on Tuesday night.",
			want:  30,
		},
		"text contains a non-encoded char": {
			given: "I love 99% of all cakes.",
			want:  8,
		},
	}

	client := NewHttpClient([]string{"tokenize", "ssplit", "pos", "lemma", "parse", "depparse"})
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pb := &nlp.Document{}
			err := client.RunText(context.Background(), []byte(tt.given), pb)
			if err != nil {
				t.Fatal(err)
			}

			if pb.Text == nil {
				t.Error("protobuf text should not be nil")
			}
			if got := tokenCount(pb); got != tt.want {
				t.Errorf("got tokens %d, want %d", got, tt.want)
			}
		})
	}
}

func tokenCount(d *nlp.Document) int {
	count := 0
	for _, s := range d.Sentence {
		count += len(s.Token)
	}
	return count
}
