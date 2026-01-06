package client

import (
	"reflect"
	"testing"

	"github.com/genelet/corenlp-golang/nlp"
)

func TestExtractTokens(t *testing.T) {
	tests := []struct {
		name string
		doc  *nlp.Document
		want []string
	}{
		{
			name: "nil document",
			doc:  nil,
			want: nil,
		},
		{
			name: "empty document",
			doc:  &nlp.Document{},
			want: nil,
		},
		{
			name: "document with tokens",
			doc: &nlp.Document{
				Sentence: []*nlp.Sentence{
					{
						Token: []*nlp.Token{
							{Word: stringPtr("Hello")},
							{Word: stringPtr("world")},
						},
					},
					{
						Token: []*nlp.Token{
							{Word: stringPtr("Test")},
						},
					},
				},
			},
			want: []string{"Hello", "world", "Test"},
		},
		{
			name: "tokens with nil words",
			doc: &nlp.Document{
				Sentence: []*nlp.Sentence{
					{
						Token: []*nlp.Token{
							{Word: stringPtr("Hello")},
							{Word: nil},
							{Word: stringPtr("world")},
						},
					},
				},
			},
			want: []string{"Hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTokens(tt.doc)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractLemmas(t *testing.T) {
	doc := &nlp.Document{
		Sentence: []*nlp.Sentence{
			{
				Token: []*nlp.Token{
					{Lemma: stringPtr("be")},
					{Lemma: stringPtr("test")},
					{Lemma: nil},
				},
			},
		},
	}

	want := []string{"be", "test"}
	got := ExtractLemmas(doc)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractLemmas() = %v, want %v", got, want)
	}
}

func TestExtractPOSTags(t *testing.T) {
	doc := &nlp.Document{
		Sentence: []*nlp.Sentence{
			{
				Token: []*nlp.Token{
					{Pos: stringPtr("NNP")},
					{Pos: stringPtr("VBZ")},
					{Pos: nil},
				},
			},
		},
	}

	want := []string{"NNP", "VBZ"}
	got := ExtractPOSTags(doc)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractPOSTags() = %v, want %v", got, want)
	}
}

func TestExtractSentences(t *testing.T) {
	doc := &nlp.Document{
		Sentence: []*nlp.Sentence{
			{
				Token: []*nlp.Token{
					{Word: stringPtr("Hello")},
					{Word: stringPtr("world")},
					{Word: stringPtr(".")},
				},
			},
			{
				Token: []*nlp.Token{
					{Word: stringPtr("Test")},
					{Word: stringPtr("sentence")},
				},
			},
		},
	}

	want := []string{"Hello world .", "Test sentence"}
	got := ExtractSentences(doc)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractSentences() = %v, want %v", got, want)
	}
}

func TestExtractNamedEntities(t *testing.T) {
	doc := &nlp.Document{
		Sentence: []*nlp.Sentence{
			{
				Token: []*nlp.Token{
					{Word: stringPtr("Stanford"), Ner: stringPtr("ORGANIZATION")},
					{Word: stringPtr("University"), Ner: stringPtr("ORGANIZATION")},
					{Word: stringPtr("is"), Ner: stringPtr("O")},
					{Word: stringPtr("in"), Ner: stringPtr("O")},
					{Word: stringPtr("California"), Ner: stringPtr("LOCATION")},
					{Word: stringPtr("."), Ner: stringPtr("O")},
				},
			},
			{
				Token: []*nlp.Token{
					{Word: stringPtr("John"), Ner: stringPtr("PERSON")},
					{Word: stringPtr("Smith"), Ner: stringPtr("PERSON")},
					{Word: stringPtr("works"), Ner: stringPtr("O")},
					{Word: stringPtr("at"), Ner: stringPtr("O")},
					{Word: stringPtr("Google"), Ner: stringPtr("ORGANIZATION")},
				},
			},
		},
	}

	want := map[string][]string{
		"ORGANIZATION": {"Stanford University", "Google"},
		"LOCATION":     {"California"},
		"PERSON":       {"John Smith"},
	}

	got := ExtractNamedEntities(doc)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractNamedEntities() = %v, want %v", got, want)
	}
}

func TestExtractTokensWithMetadata(t *testing.T) {
	doc := &nlp.Document{
		Sentence: []*nlp.Sentence{
			{
				Token: []*nlp.Token{
					{
						Word:  stringPtr("Stanford"),
						Lemma: stringPtr("Stanford"),
						Pos:   stringPtr("NNP"),
						Ner:   stringPtr("ORGANIZATION"),
					},
					{
						Word:  stringPtr("is"),
						Lemma: stringPtr("be"),
						Pos:   stringPtr("VBZ"),
						Ner:   stringPtr("O"),
					},
				},
			},
		},
	}

	want := []TokenWithMetadata{
		{Word: "Stanford", Lemma: "Stanford", POS: "NNP", NER: "ORGANIZATION"},
		{Word: "is", Lemma: "be", POS: "VBZ", NER: "O"},
	}

	got := ExtractTokensWithMetadata(doc)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractTokensWithMetadata() = %v, want %v", got, want)
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
