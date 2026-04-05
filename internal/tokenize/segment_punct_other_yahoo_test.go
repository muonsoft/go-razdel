package tokenize_test

import (
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// T008: punct / other / yahoo join rules mirror
// third_party/razdel/razdel/segmenters/tokenize.py (FunctionRule punct, other, yahoo)
// and razdel/tests/test_tokenize.py UNIT cases.

func TestTokenTexts_punct_ellipsisAndEndingRuns(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"...":   {"..."},
		"...?!": {"...?!"},
		"?!.":   {"?!."},
		"!?.":   {"!?."},
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

func TestTokenTexts_punct_smileColonParens(t *testing.T) {
	t.Parallel()
	// SMILE in punct.py matches [=:;]-?[)(]{1,3} on the accumulated buffer + next atom.
	got := tokenize.TokenTexts(":)))")
	want := []string{":)))"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	got = tokenize.TokenTexts(":)||,")
	want = []string{":)", "|", "|", ","}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestTokenTexts_punct_starRun(t *testing.T) {
	t.Parallel()
	got := tokenize.TokenTexts("***")
	want := []string{"***"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestTokenTexts_other_mixedScript(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"mβж": {"mβж"},
		"Δσ":  {"Δσ"},
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

func TestTokenTexts_yahoo_exclamation(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"yahoo!": {"yahoo!"},
		"Yahoo!": {"Yahoo!"},
		"YAHOO!": {"YAHOO!"},
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

func TestTokenTexts_punct_quoteBracketPipeUnit(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"»||.": {"»", "|", "|", "."},
		")||.": {")", "|", "|", "."},
		"(||«": {"(", "|", "|", "«"},
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

func TestTokenTexts_T008_regression_dashFloatFractionUnchanged(t *testing.T) {
	t.Parallel()
	cases := map[string][]string{
		"что-то": {"что-то"},
		"1,5":    {"1,5"},
		"1.5":    {"1.5"},
		"1/2":    {"1/2"},
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
