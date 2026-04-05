// Package razdel provides Russian text tokenization and sentence segmentation.
//
// Public API and offset semantics are defined in docs/contracts.md (UTF-8 byte offsets, Variant A).
package razdel

import "github.com/muonsoft/go-razdel/internal/tokenize"

// Span is a half-open byte interval into the original UTF-8 string: [Start, End).
// Start and End are measured in bytes; see docs/contracts.md.
type Span struct {
	Start int
	End   int
}

// Token is a token span with its source text slice.
type Token struct {
	Span
	Text string
}

// Sentence is a sentence span with its source text slice.
type Sentence struct {
	Span
	Text string
}

// Tokenize splits text into tokens with behavior aligned to upstream
// third_party/razdel/razdel/segmenters/tokenize.py (atoms, splits, join rules).
func Tokenize(text string) []Token {
	if text == "" {
		return nil
	}
	spans := tokenize.TokenSpans(text)
	if len(spans) == 0 {
		return nil
	}
	out := make([]Token, len(spans))
	for i, s := range spans {
		start, end := s[0], s[1]
		out[i] = Token{
			Span: Span{Start: start, End: end},
			Text: text[start:end],
		}
	}
	return out
}

// Sentenize splits text into sentences. Stub: returns an empty slice for any input.
// Split stream and SentSplit accessors live in internal/sentenize (upstream SentSplitter / SentSplit).
func Sentenize(text string) []Sentence {
	_ = text
	return nil
}
