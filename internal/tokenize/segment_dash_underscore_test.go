package tokenize_test

import (
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// T006: DashRule / UnderscoreRule (Rule2112) — join when the dash or underscore is
// the adjacent atom text and neither neighbor across the join is PUNCT.
// Mirrors third_party/razdel/razdel/segmenters/tokenize.py (DashRule, UnderscoreRule).

func TestTokenTexts_joinDash_cyrillicCompound(t *testing.T) {
	t.Parallel()
	got := tokenize.TokenTexts("что-то")
	want := []string{"что-то"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestTokenTexts_joinUnderscore_cyrillicPhrase(t *testing.T) {
	t.Parallel()
	got := tokenize.TokenTexts("К_тому_же")
	want := []string{"К_тому_же"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestTokenTexts_noJoin_dashBetweenPunctuationAtoms(t *testing.T) {
	t.Parallel()
	// Upstream: PUNCT–dash–PUNCT does not join (!-! → '!', '-', '!').
	got := tokenize.TokenTexts("!-!")
	want := []string{"!", "-", "!"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestTokenTexts_noJoin_underscoreBetweenPunctuationAtoms(t *testing.T) {
	t.Parallel()
	got := tokenize.TokenTexts("!_!")
	want := []string{"!", "_", "!"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
