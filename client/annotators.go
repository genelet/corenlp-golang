package client

import (
	"fmt"
	"strings"
)

// Annotator represents a Stanford CoreNLP annotator.
// Use the predefined constants to avoid typos and get IDE autocomplete support.
type Annotator string

// Core annotators - fundamental NLP operations
const (
	// AnnotatorTokenize splits the raw text into individual tokens (words).
	// This is typically the first annotator in any pipeline.
	AnnotatorTokenize Annotator = "tokenize"

	// AnnotatorCleanXML removes XML tags from the document.
	// This should typically run before tokenization.
	AnnotatorCleanXML Annotator = "cleanxml"

	// AnnotatorSSplit divides tokenized text into sentences.
	// Requires: tokenize
	AnnotatorSSplit Annotator = "ssplit"

	// AnnotatorDocDate extracts the document date.
	// This can be used independently or as a dependency for other annotators.
	AnnotatorDocDate Annotator = "docdate"

	// AnnotatorPOS assigns part-of-speech tags to tokens (e.g., noun, verb, adjective).
	// Requires: tokenize, ssplit
	AnnotatorPOS Annotator = "pos"

	// AnnotatorLemma determines the base or dictionary form (lemma) for each word.
	// Requires: tokenize, ssplit, pos
	AnnotatorLemma Annotator = "lemma"
)

// Named Entity Recognition annotators
const (
	// AnnotatorNER identifies and classifies named entities (persons, organizations, locations, dates).
	// Requires: tokenize, ssplit, pos, lemma
	AnnotatorNER Annotator = "ner"

	// AnnotatorRegexNER implements rule-based NER using Java regular expressions.
	// Can identify entities not in traditional NLP corpora (ideologies, nationalities, religions, titles).
	// Requires: tokenize, ssplit, pos, lemma, ner
	AnnotatorRegexNER Annotator = "regexner"

	// AnnotatorEntityMentions detects full entity mentions.
	// Often run automatically as a sub-annotator of NER.
	// Requires: tokenize, ssplit, pos, lemma, ner
	AnnotatorEntityMentions Annotator = "entitymentions"

	// AnnotatorEntityLink links entity mentions to Wikipedia entities.
	// Requires: tokenize, ssplit, pos, lemma, ner, entitymentions
	AnnotatorEntityLink Annotator = "entitylink"
)

// Parsing annotators - syntactic analysis
const (
	// AnnotatorParse generates a full syntactic parse tree (constituency parsing).
	// Requires: tokenize, ssplit, pos
	AnnotatorParse Annotator = "parse"

	// AnnotatorDepparse provides dependency parsing with basic and enhanced dependencies.
	// Requires: tokenize, ssplit, pos
	AnnotatorDepparse Annotator = "depparse"
)

// Coreference resolution annotators
const (
	// AnnotatorCoref identifies and resolves coreferences (linking mentions to the same entity).
	// Requires: tokenize, ssplit, pos, lemma, ner, parse
	AnnotatorCoref Annotator = "coref"

	// AnnotatorDcoref is the older deterministic coreference resolution system.
	// Deprecated in favor of AnnotatorCoref.
	// Requires: tokenize, ssplit, pos, lemma, ner, parse
	AnnotatorDcoref Annotator = "dcoref"

	// AnnotatorMention identifies coreference mentions.
	// Often run automatically as a sub-annotator of coref.
	// Requires: tokenize, ssplit, pos, lemma, ner, parse
	AnnotatorMention Annotator = "mention"
)

// Sentiment and semantic annotators
const (
	// AnnotatorSentiment determines the sentiment of sentences.
	// Requires: tokenize, ssplit, parse
	AnnotatorSentiment Annotator = "sentiment"

	// AnnotatorNatlog marks quantifier scope and token polarity for natural logic semantics.
	// Requires: tokenize, ssplit, pos, lemma, depparse
	AnnotatorNatlog Annotator = "natlog"

	// AnnotatorOpenie extracts open-domain relation triples from sentences.
	// Requires: tokenize, ssplit, pos, lemma, depparse, natlog
	AnnotatorOpenie Annotator = "openie"

	// AnnotatorTruecase determines the true case of tokens in the text.
	// Useful for text that is all lowercase or all uppercase.
	AnnotatorTruecase Annotator = "truecase"

	// AnnotatorUDFeats adds Universal Dependencies features to tokens.
	// Requires: tokenize, ssplit, pos
	AnnotatorUDFeats Annotator = "udfeats"
)

// Information extraction annotators
const (
	// AnnotatorRelation finds relations between two entities based on a trained model.
	// Default relations: Live_In, Located_In, OrgBased_In, Work_For, None.
	// Requires: tokenize, ssplit, pos, lemma, ner, parse
	AnnotatorRelation Annotator = "relation"

	// AnnotatorKBP extracts (subject, relation, object) triples for Knowledge Base Population.
	// Requires: tokenize, ssplit, pos, lemma, ner, parse, coref
	AnnotatorKBP Annotator = "kbp"
)

// Text and quote processing
const (
	// AnnotatorQuote identifies quotes delimited by single or double quotation marks.
	// Requires: tokenize, ssplit
	AnnotatorQuote Annotator = "quote"

	// AnnotatorQuoteAttribution attributes quotes to speakers.
	// Requires: tokenize, ssplit, pos, lemma, ner, depparse, coref, quote
	AnnotatorQuoteAttribution Annotator = "quote.attribution"

	// AnnotatorTokensRegex runs TokensRegex patterns within the NLP pipeline.
	// Requires: tokenize, ssplit
	AnnotatorTokensRegex Annotator = "tokensregex"
)

// Predefined annotator combinations for common use cases.
// These provide sensible defaults for typical NLP tasks.
var (
	// BasicAnnotators provides tokenization, sentence splitting, POS tagging, and lemmatization.
	// This is the foundation for most NLP pipelines.
	BasicAnnotators = []Annotator{
		AnnotatorTokenize,
		AnnotatorSSplit,
		AnnotatorPOS,
		AnnotatorLemma,
	}

	// SyntaxAnnotators adds constituency and dependency parsing to basic annotators.
	// Use this for syntactic analysis tasks.
	SyntaxAnnotators = []Annotator{
		AnnotatorTokenize,
		AnnotatorSSplit,
		AnnotatorPOS,
		AnnotatorLemma,
		AnnotatorParse,
		AnnotatorDepparse,
	}

	// NERAnnotators adds named entity recognition to basic annotators.
	// Use this for entity extraction tasks.
	NERAnnotators = []Annotator{
		AnnotatorTokenize,
		AnnotatorSSplit,
		AnnotatorPOS,
		AnnotatorLemma,
		AnnotatorNER,
		AnnotatorEntityMentions,
	}

	// SemanticAnnotators provides a comprehensive pipeline including parsing, NER, and coreference.
	// Use this for deep semantic analysis.
	SemanticAnnotators = []Annotator{
		AnnotatorTokenize,
		AnnotatorSSplit,
		AnnotatorPOS,
		AnnotatorLemma,
		AnnotatorNER,
		AnnotatorParse,
		AnnotatorDepparse,
		AnnotatorCoref,
	}

	// RelationExtractionAnnotators provides everything needed for relation extraction.
	RelationExtractionAnnotators = []Annotator{
		AnnotatorTokenize,
		AnnotatorSSplit,
		AnnotatorPOS,
		AnnotatorLemma,
		AnnotatorNER,
		AnnotatorParse,
		AnnotatorDepparse,
		AnnotatorNatlog,
		AnnotatorOpenie,
	}
)

// String returns the annotator as a string.
func (a Annotator) String() string {
	return string(a)
}

// ValidateAnnotators checks if all provided annotators are valid.
// Returns an error if any annotator is empty or contains only whitespace.
// Note: This does not check for annotator dependencies or ordering.
func ValidateAnnotators(annotators []Annotator) error {
	if len(annotators) == 0 {
		return fmt.Errorf("at least one annotator is required")
	}

	for i, ann := range annotators {
		if strings.TrimSpace(string(ann)) == "" {
			return fmt.Errorf("annotator at index %d is empty or whitespace", i)
		}
	}

	return nil
}

// AnnotatorsToStrings converts a slice of Annotator to a slice of strings.
// This is useful for compatibility with functions that expect []string.
func AnnotatorsToStrings(annotators []Annotator) []string {
	result := make([]string, len(annotators))
	for i, ann := range annotators {
		result[i] = string(ann)
	}
	return result
}

// StringsToAnnotators converts a slice of strings to a slice of Annotator.
// This is useful for converting existing string-based annotator lists.
func StringsToAnnotators(annotators []string) []Annotator {
	result := make([]Annotator, len(annotators))
	for i, ann := range annotators {
		result[i] = Annotator(ann)
	}
	return result
}
