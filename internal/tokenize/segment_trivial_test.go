package tokenize_test

import (
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// T005: trivial split-space baseline — tokens split at non-empty delimiters
// (upstream split_space + TokenSegmenter.segment), independent of join rules.
func TestTokenTexts_trivialSplit_space(t *testing.T) {
	t.Parallel()
	got := tokenize.TokenTexts("a b")
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	spans := tokenize.TokenSpans("a b")
	if !reflect.DeepEqual(spans, [][2]int{{0, 1}, {2, 3}}) {
		t.Fatalf("TokenSpans: got %#v", spans)
	}
}

func TestTokenTexts_trivialSplit_punctuationThenSpace(t *testing.T) {
	t.Parallel()
	src := "привет, мир"
	// Comma is its own atom; space before "мир" is delimiter => split (T005 doc).
	got := tokenize.TokenTexts(src)
	want := []string{"привет", ",", "мир"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	n0 := len("привет")
	nComma := len(",")
	nSpace := len(" ")
	n2 := len("мир")
	wantSpans := [][2]int{
		{0, n0},
		{n0, n0 + nComma},
		{n0 + nComma + nSpace, n0 + nComma + nSpace + n2},
	}
	spans := tokenize.TokenSpans(src)
	if !reflect.DeepEqual(spans, wantSpans) {
		t.Fatalf("TokenSpans: got %#v, want %#v", spans, wantSpans)
	}
}
