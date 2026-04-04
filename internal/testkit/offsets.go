package testkit

import (
	"fmt"
	"testing"

	"github.com/muonsoft/go-razdel"
)

// ValidateTokenOffsets checks Variant A offsets for each token: valid range and Text == src[Start:End].
func ValidateTokenOffsets(src string, tokens []razdel.Token) error {
	for i, tok := range tokens {
		if err := validateSpanText("Token", i, src, tok.Start, tok.End, tok.Text); err != nil {
			return err
		}
	}
	return nil
}

// ValidateSentenceOffsets checks Variant A offsets for each sentence.
func ValidateSentenceOffsets(src string, sents []razdel.Sentence) error {
	for i, s := range sents {
		if err := validateSpanText("Sentence", i, src, s.Start, s.End, s.Text); err != nil {
			return err
		}
	}
	return nil
}

// AssertTokenOffsetContract checks Variant A offsets for each token: valid range and Text == src[Start:End].
func AssertTokenOffsetContract(tb testing.TB, src string, tokens []razdel.Token) {
	tb.Helper()
	if err := ValidateTokenOffsets(src, tokens); err != nil {
		tb.Fatal(err)
	}
}

// AssertSentenceOffsetContract checks Variant A offsets for each sentence.
func AssertSentenceOffsetContract(tb testing.TB, src string, sents []razdel.Sentence) {
	tb.Helper()
	if err := ValidateSentenceOffsets(src, sents); err != nil {
		tb.Fatal(err)
	}
}

// AssertTokenTextsEqual checks len(tokens) == len(want) and each token's Text matches want[i]
// after AssertTokenOffsetContract passes.
func AssertTokenTextsEqual(tb testing.TB, src string, tokens []razdel.Token, want []string) {
	tb.Helper()
	AssertTokenOffsetContract(tb, src, tokens)
	if len(tokens) != len(want) {
		tb.Fatalf("token count: want %d, got %d", len(want), len(tokens))
	}
	for i := range tokens {
		if tokens[i].Text != want[i] {
			tb.Fatalf("token[%d].Text: want %q, got %q", i, want[i], tokens[i].Text)
		}
	}
}

// AssertSentenceTextsEqual checks sentence texts against want.
func AssertSentenceTextsEqual(tb testing.TB, src string, sents []razdel.Sentence, want []string) {
	tb.Helper()
	AssertSentenceOffsetContract(tb, src, sents)
	if len(sents) != len(want) {
		tb.Fatalf("sentence count: want %d, got %d", len(want), len(sents))
	}
	for i := range sents {
		if sents[i].Text != want[i] {
			tb.Fatalf("sentence[%d].Text: want %q, got %q", i, want[i], sents[i].Text)
		}
	}
}

func validateSpanText(kind string, idx int, src string, start, end int, text string) error {
	n := len(src)
	if err := ValidateHalfOpen(start, end, n); err != nil {
		return fmt.Errorf("%s[%d]: %w", kind, idx, err)
	}
	slice := src[start:end]
	if slice != text {
		return fmt.Errorf("%s[%d]: Text %q != src[%d:%d] %q", kind, idx, text, start, end, slice)
	}
	return nil
}
