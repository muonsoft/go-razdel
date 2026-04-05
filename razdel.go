// Package razdel provides Russian text tokenization and sentence segmentation.
//
// Public API and offset semantics are defined in docs/contracts.md (UTF-8 byte offsets, Variant A).
package razdel

import (
	"strings"

	"github.com/muonsoft/go-razdel/internal/sentenize"
	"github.com/muonsoft/go-razdel/internal/tokenize"
)

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

// Sentenize splits text into sentences using the trivial join layer from upstream
// sentenize.py (empty_side, no_space_prefix, lower_right, delimiter_right). Further
// upstream rules (sokr, bounds, bullets, dash) are not applied yet; parity grows with later tasks.
func Sentenize(text string) []Sentence {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	parts := sentenize.SentSplitterParts(text, sentenize.DefaultWindow)
	raw := sentenize.Segment(parts, sentenize.JoinTrivial)
	chunks := sentenize.PostStrip(raw)
	spans := sentenize.ByteSpans(text, chunks)
	out := make([]Sentence, 0, len(spans))
	for _, sp := range spans {
		start, end := sp[0], sp[1]
		out = append(out, Sentence{
			Span: Span{Start: start, End: end},
			Text: text[start:end],
		})
	}
	return out
}
