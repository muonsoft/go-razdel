package razdel_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

// Synthetic Token/Sentence slices exercise Variant A offsets (docs/contracts.md) independent of tokenizer logic.
func TestContract_tokenOffsets_table(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		src    string
		tokens []razdel.Token
	}{
		{
			name: "ascii_partition",
			src:  "a b",
			tokens: []razdel.Token{
				{Span: razdel.Span{Start: 0, End: 1}, Text: "a"},
				{Span: razdel.Span{Start: 1, End: 2}, Text: " "},
				{Span: razdel.Span{Start: 2, End: 3}, Text: "b"},
			},
		},
		{
			name: "cyrillic_multibyte",
			src:  "ж",
			tokens: []razdel.Token{
				{Span: razdel.Span{Start: 0, End: len("ж")}, Text: "ж"},
			},
		},
		{
			name: "emoji_multibyte",
			src:  "🙂",
			tokens: []razdel.Token{
				{Span: razdel.Span{Start: 0, End: len("🙂")}, Text: "🙂"},
			},
		},
		{
			name: "mixed_ascii_cyrillic_emoji",
			src:  "aж🙂",
			tokens: []razdel.Token{
				{Span: razdel.Span{Start: 0, End: 1}, Text: "a"},
				{Span: razdel.Span{Start: 1, End: 1 + len("ж")}, Text: "ж"},
				{Span: razdel.Span{Start: 1 + len("ж"), End: len("aж🙂")}, Text: "🙂"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testkit.AssertTokenOffsetContract(t, tt.src, tt.tokens)
			for i, tok := range tt.tokens {
				if err := testkit.ValidateHalfOpen(tok.Start, tok.End, len(tt.src)); err != nil {
					t.Fatalf("token[%d] half-open: %v", i, err)
				}
			}
		})
	}
}

func TestContract_sentenceOffsets_table(t *testing.T) {
	t.Parallel()
	src := "Первое. Второе."
	tests := []struct {
		name string
		s    []razdel.Sentence
	}{
		{
			name: "two_sentences",
			s: []razdel.Sentence{
				{Span: razdel.Span{Start: 0, End: len("Первое.")}, Text: "Первое."},
				{Span: razdel.Span{Start: len("Первое.") + 1, End: len(src)}, Text: "Второе."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testkit.AssertSentenceOffsetContract(t, src, tt.s)
		})
	}
}

func TestContract_tokenizeSentenize_deterministic(t *testing.T) {
	t.Parallel()
	inputs := []string{
		"",
		"a",
		"привет",
		"aж🙂",
		"Line one.\nLine two.",
	}
	for _, s := range inputs {
		s := s
		t.Run(fmt.Sprintf("%q", s), func(t *testing.T) {
			t.Parallel()
			tok1 := razdel.Tokenize(s)
			tok2 := razdel.Tokenize(s)
			if !reflect.DeepEqual(tok1, tok2) {
				t.Fatalf("Tokenize not deterministic for %q", s)
			}
			sent1 := razdel.Sentenize(s)
			sent2 := razdel.Sentenize(s)
			if !reflect.DeepEqual(sent1, sent2) {
				t.Fatalf("Sentenize not deterministic for %q", s)
			}
		})
	}
}

func TestContract_tokenizeSentenize_validUTF8_noPanic(t *testing.T) {
	t.Parallel()
	cases := []string{
		"",
		"ascii",
		"кириллица",
		"🙂🔥",
		"mixed привет 🙂",
	}
	for _, s := range cases {
		s := s
		t.Run(fmt.Sprintf("%q", s), func(t *testing.T) {
			t.Parallel()
			_ = razdel.Tokenize(s)
			_ = razdel.Sentenize(s)
		})
	}
}
