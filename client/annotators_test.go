package client

import (
	"testing"
)

func TestAnnotatorConstants(t *testing.T) {
	// Test that annotator constants return correct strings
	tests := []struct {
		annotator Annotator
		expected  string
	}{
		{AnnotatorTokenize, "tokenize"},
		{AnnotatorSSplit, "ssplit"},
		{AnnotatorPOS, "pos"},
		{AnnotatorLemma, "lemma"},
		{AnnotatorNER, "ner"},
		{AnnotatorParse, "parse"},
		{AnnotatorDepparse, "depparse"},
		{AnnotatorCoref, "coref"},
		{AnnotatorSentiment, "sentiment"},
		{AnnotatorOpenie, "openie"},
	}

	for _, tt := range tests {
		t.Run(string(tt.annotator), func(t *testing.T) {
			if got := tt.annotator.String(); got != tt.expected {
				t.Errorf("Annotator.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateAnnotators(t *testing.T) {
	tests := []struct {
		name        string
		annotators  []Annotator
		expectError bool
	}{
		{
			name:        "valid annotators",
			annotators:  []Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS},
			expectError: false,
		},
		{
			name:        "empty slice",
			annotators:  []Annotator{},
			expectError: true,
		},
		{
			name:        "nil slice",
			annotators:  nil,
			expectError: true,
		},
		{
			name:        "contains empty string",
			annotators:  []Annotator{AnnotatorTokenize, "", AnnotatorPOS},
			expectError: true,
		},
		{
			name:        "contains whitespace",
			annotators:  []Annotator{AnnotatorTokenize, "  ", AnnotatorPOS},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAnnotators(tt.annotators)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateAnnotators() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestAnnotatorsToStrings(t *testing.T) {
	annotators := []Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS}
	expected := []string{"tokenize", "ssplit", "pos"}

	result := AnnotatorsToStrings(annotators)

	if len(result) != len(expected) {
		t.Fatalf("AnnotatorsToStrings() length = %v, want %v", len(result), len(expected))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("AnnotatorsToStrings()[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestStringsToAnnotators(t *testing.T) {
	strings := []string{"tokenize", "ssplit", "pos"}
	expected := []Annotator{AnnotatorTokenize, AnnotatorSSplit, AnnotatorPOS}

	result := StringsToAnnotators(strings)

	if len(result) != len(expected) {
		t.Fatalf("StringsToAnnotators() length = %v, want %v", len(result), len(expected))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("StringsToAnnotators()[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestPredefinedAnnotatorCombinations(t *testing.T) {
	tests := []struct {
		name        string
		combination []Annotator
		minLength   int
	}{
		{"BasicAnnotators", BasicAnnotators, 4},
		{"SyntaxAnnotators", SyntaxAnnotators, 6},
		{"NERAnnotators", NERAnnotators, 6},
		{"SemanticAnnotators", SemanticAnnotators, 8},
		{"RelationExtractionAnnotators", RelationExtractionAnnotators, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.combination) < tt.minLength {
				t.Errorf("%s length = %v, want at least %v", tt.name, len(tt.combination), tt.minLength)
			}

			// Verify no empty annotators
			for i, ann := range tt.combination {
				if ann == "" {
					t.Errorf("%s[%d] is empty", tt.name, i)
				}
			}

			// Verify all combinations include tokenize and ssplit (foundation)
			if tt.combination[0] != AnnotatorTokenize {
				t.Errorf("%s should start with tokenize", tt.name)
			}
			if tt.combination[1] != AnnotatorSSplit {
				t.Errorf("%s should have ssplit as second annotator", tt.name)
			}
		})
	}
}
