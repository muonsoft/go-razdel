package tokenize_test

import (
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// IGO-47: smile lookahead (upstream #17) and emoji/word split (upstream #2).
// Go intentionally drifts from pinned Python tokenize for these cases; see README.md.

func TestTokenTexts_IGO47_smileLookahead(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		":-)":             {":-)"},
		";-)":             {";-)"},
		"=-)":             {"=-)"},
		":)))":            {":)))"},
		"text :-) text":   {"text", ":-)", "text"},
		"hello ;-) there": {"hello", ";-)", "there"},
	}
	for input, want := range cases {
		input, want := input, want
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			got := tokenize.TokenTexts(input)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v, want %#v", got, want)
			}
		})
	}
}

func TestTokenTexts_IGO47_emojiSplitFromWords(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"✅Сдается":   {"✅", "Сдается"},
		"счетчики💰": {"счетчики", "💰"},
		"abc📲def":   {"abc", "📲", "def"},
		"📲":         {"📲"},
	}
	for input, want := range cases {
		input, want := input, want
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			got := tokenize.TokenTexts(input)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v, want %#v", got, want)
			}
		})
	}
}

func TestTokenSpans_IGO47_offsets(t *testing.T) {
	t.Parallel()
	src := "ab✅cd"
	spans := tokenize.TokenSpans(src)
	if len(spans) != 3 {
		t.Fatalf("len %d, want 3", len(spans))
	}
	for i, want := range []string{"ab", "✅", "cd"} {
		s, e := spans[i][0], spans[i][1]
		if got := src[s:e]; got != want {
			t.Fatalf("token %d: got %q want %q", i, got, want)
		}
	}
}
