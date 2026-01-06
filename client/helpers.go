package client

import (
	"strings"

	"github.com/genelet/corenlp-golang/nlp"
)

// ExtractTokens extracts all token words from a document as a slice of strings.
// This convenience method iterates through all sentences and collects token words.
//
// Example:
//
//	doc := &nlp.Document{}
//	client.RunText(ctx, text, doc)
//	tokens := ExtractTokens(doc)
//	// tokens = ["Stanford", "University", "is", "located", "in", "California", "."]
func ExtractTokens(doc *nlp.Document) []string {
	if doc == nil {
		return nil
	}

	var tokens []string
	for _, sentence := range doc.Sentence {
		for _, token := range sentence.Token {
			if token.Word != nil {
				tokens = append(tokens, *token.Word)
			}
		}
	}
	return tokens
}

// ExtractSentences extracts all sentences from a document as a slice of strings.
// Each sentence is reconstructed from its tokens' words.
//
// Example:
//
//	doc := &nlp.Document{}
//	client.RunText(ctx, text, doc)
//	sentences := ExtractSentences(doc)
//	// sentences = ["Stanford University is located in California ."]
func ExtractSentences(doc *nlp.Document) []string {
	if doc == nil {
		return nil
	}

	sentences := make([]string, 0, len(doc.Sentence))
	for _, sentence := range doc.Sentence {
		var words []string
		for _, token := range sentence.Token {
			if token.Word != nil {
				words = append(words, *token.Word)
			}
		}
		if len(words) > 0 {
			sentences = append(sentences, strings.Join(words, " "))
		}
	}
	return sentences
}

// ExtractLemmas extracts all lemmas from a document as a slice of strings.
// Lemmas are the base or dictionary forms of words.
//
// Example:
//
//	doc := &nlp.Document{}
//	client.RunText(ctx, text, doc)
//	lemmas := ExtractLemmas(doc)
//	// lemmas = ["Stanford", "University", "be", "locate", "in", "California", "."]
func ExtractLemmas(doc *nlp.Document) []string {
	if doc == nil {
		return nil
	}

	var lemmas []string
	for _, sentence := range doc.Sentence {
		for _, token := range sentence.Token {
			if token.Lemma != nil {
				lemmas = append(lemmas, *token.Lemma)
			}
		}
	}
	return lemmas
}

// ExtractPOSTags extracts all part-of-speech tags from a document as a slice of strings.
//
// Example:
//
//	doc := &nlp.Document{}
//	client.RunText(ctx, text, doc)
//	posTags := ExtractPOSTags(doc)
//	// posTags = ["NNP", "NNP", "VBZ", "VBN", "IN", "NNP", "."]
func ExtractPOSTags(doc *nlp.Document) []string {
	if doc == nil {
		return nil
	}

	var tags []string
	for _, sentence := range doc.Sentence {
		for _, token := range sentence.Token {
			if token.Pos != nil {
				tags = append(tags, *token.Pos)
			}
		}
	}
	return tags
}

// ExtractNamedEntities extracts all named entities from a document grouped by entity type.
// Returns a map where keys are NER tags (e.g., "PERSON", "LOCATION") and values are slices of entity strings.
//
// Example:
//
//	doc := &nlp.Document{}
//	client.RunText(ctx, text, doc)
//	entities := ExtractNamedEntities(doc)
//	// entities = {"ORGANIZATION": ["Stanford University"], "LOCATION": ["California"]}
func ExtractNamedEntities(doc *nlp.Document) map[string][]string {
	if doc == nil {
		return nil
	}

	entities := make(map[string][]string)

	for _, sentence := range doc.Sentence {
		var currentNER string
		var currentEntity []string

		for _, token := range sentence.Token {
			ner := ""
			if token.Ner != nil {
				ner = *token.Ner
			}

			// If NER tag changes or is "O" (outside entity), save current entity
			if ner != currentNER || ner == "O" || ner == "" {
				if currentNER != "" && currentNER != "O" && len(currentEntity) > 0 {
					entityText := strings.Join(currentEntity, " ")
					entities[currentNER] = append(entities[currentNER], entityText)
				}
				currentEntity = nil
				currentNER = ner
			}

			// Add to current entity if not "O"
			if ner != "" && ner != "O" && token.Word != nil {
				currentEntity = append(currentEntity, *token.Word)
			}
		}

		// Save any remaining entity at end of sentence
		if currentNER != "" && currentNER != "O" && len(currentEntity) > 0 {
			entityText := strings.Join(currentEntity, " ")
			entities[currentNER] = append(entities[currentNER], entityText)
		}
	}

	return entities
}

// TokenWithMetadata represents a token with its common metadata fields.
type TokenWithMetadata struct {
	Word  string
	Lemma string
	POS   string
	NER   string
}

// ExtractTokensWithMetadata extracts tokens along with their lemma, POS, and NER information.
// This is useful when you need multiple fields from each token.
//
// Example:
//
//	doc := &nlp.Document{}
//	client.RunText(ctx, text, doc)
//	tokens := ExtractTokensWithMetadata(doc)
//	for _, token := range tokens {
//	    fmt.Printf("%s (%s, %s)\n", token.Word, token.POS, token.Lemma)
//	}
func ExtractTokensWithMetadata(doc *nlp.Document) []TokenWithMetadata {
	if doc == nil {
		return nil
	}

	var tokens []TokenWithMetadata
	for _, sentence := range doc.Sentence {
		for _, token := range sentence.Token {
			t := TokenWithMetadata{}
			if token.Word != nil {
				t.Word = *token.Word
			}
			if token.Lemma != nil {
				t.Lemma = *token.Lemma
			}
			if token.Pos != nil {
				t.POS = *token.Pos
			}
			if token.Ner != nil {
				t.NER = *token.Ner
			}
			tokens = append(tokens, t)
		}
	}
	return tokens
}
