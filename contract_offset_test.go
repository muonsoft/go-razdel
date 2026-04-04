package razdel_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

// Synthetic Token/Sentence slices exercise Variant A offsets (docs/contracts.md) via testkit.
// They do not call Tokenize/Sentenize; extend with real API assertions when T004+ implements segmentation.
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
		})
	}
}

func TestContract_sentenceOffsets_table(t *testing.T) {
	t.Parallel()
	srcTwo := "Первое. Второе."
	emojiSrc := "Привет 🙂."
	tests := []struct {
		name string
		src  string
		s    []razdel.Sentence
	}{
		{
			name: "empty_string",
			src:  "",
			s:    nil,
		},
		{
			name: "single_sentence",
			src:  "Одно.",
			s: []razdel.Sentence{
				{Span: razdel.Span{Start: 0, End: len("Одно.")}, Text: "Одно."},
			},
		},
		{
			name: "two_sentences",
			src:  srcTwo,
			s: []razdel.Sentence{
				{Span: razdel.Span{Start: 0, End: len("Первое.")}, Text: "Первое."},
				{Span: razdel.Span{Start: len("Первое.") + 1, End: len(srcTwo)}, Text: "Второе."},
			},
		},
		{
			name: "sentence_with_emoji",
			src:  emojiSrc,
			s: []razdel.Sentence{
				{Span: razdel.Span{Start: 0, End: len(emojiSrc)}, Text: emojiSrc},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testkit.AssertSentenceOffsetContract(t, tt.src, tt.s)
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
		t.Run(fmt.Sprintf("%q", s), func(t *testing.T) {
			t.Parallel()
			_ = razdel.Tokenize(s)
			_ = razdel.Sentenize(s)
		})
	}
}
