package testkit

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/muonsoft/go-razdel"
)

func TestValidateTokenOffsets_ascii(t *testing.T) {
	t.Parallel()
	src := "hello world"
	tokens := []razdel.Token{
		{Span: razdel.Span{Start: 0, End: 5}, Text: "hello"},
		{Span: razdel.Span{Start: 6, End: 11}, Text: "world"},
	}
	if err := ValidateTokenOffsets(src, tokens); err != nil {
		t.Fatal(err)
	}
}

func TestValidateTokenOffsets_utf8(t *testing.T) {
	t.Parallel()
	src := "привет мир"
	if !utf8.ValidString(src) {
		t.Fatal("fixture must be valid UTF-8")
	}
	hi := strings.Index(src, "привет")
	if hi != 0 {
		t.Fatalf("unexpected prefix offset %d", hi)
	}
	hiEnd := hi + len("привет")
	space := hiEnd
	if src[space] != ' ' {
		t.Fatalf("expected space at byte %d", space)
	}
	mirStart := space + 1
	tokens := []razdel.Token{
		{Span: razdel.Span{Start: hi, End: hiEnd}, Text: "привет"},
		{Span: razdel.Span{Start: space, End: space + 1}, Text: " "},
		{Span: razdel.Span{Start: mirStart, End: len(src)}, Text: "мир"},
	}
	if err := ValidateTokenOffsets(src, tokens); err != nil {
		t.Fatal(err)
	}
}

func TestValidateTokenOffsets_invalidSpan(t *testing.T) {
	t.Parallel()
	src := "ab"
	tokens := []razdel.Token{
		{Span: razdel.Span{Start: 0, End: 3}, Text: "ab"},
	}
	if err := ValidateTokenOffsets(src, tokens); err == nil {
		t.Fatal("expected error for End > len(src)")
	}
}

func TestValidateTokenOffsets_textMismatch(t *testing.T) {
	t.Parallel()
	src := "ab"
	tokens := []razdel.Token{
		{Span: razdel.Span{Start: 0, End: 2}, Text: "xx"},
	}
	if err := ValidateTokenOffsets(src, tokens); err == nil {
		t.Fatal("expected error when Text != slice")
	}
}

func TestValidateSentenceOffsets(t *testing.T) {
	t.Parallel()
	src := "One. Two."
	sents := []razdel.Sentence{
		{Span: razdel.Span{Start: 0, End: 4}, Text: "One."},
		{Span: razdel.Span{Start: 5, End: 9}, Text: "Two."},
	}
	if err := ValidateSentenceOffsets(src, sents); err != nil {
		t.Fatal(err)
	}
}
