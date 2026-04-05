package tokenize_test

import (
	"reflect"
	"testing"

	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// T007: FloatRule / FractionRule (Rule2112) — join INT with delimiter in ., or /\
// when the delimiter is the adjacent atom text. Mirrors
// third_party/razdel/razdel/segmenters/tokenize.py (FloatRule, FractionRule).

func TestTokenTexts_joinFloat_commaAndDot(t *testing.T) {
	t.Parallel()
	for _, input := range []string{"1,5", "1.5"} {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			got := tokenize.TokenTexts(input)
			want := []string{input}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v, want %#v", got, want)
			}
		})
	}
}

func TestTokenTexts_joinFraction_slashAndBackslash(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input string
		want  []string
	}{
		{"1/2", []string{"1/2"}},
		{`1\2`, []string{`1\2`}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got := tokenize.TokenTexts(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestTokenTexts_noJoin_commaBetweenLetters(t *testing.T) {
	t.Parallel()
	got := tokenize.TokenTexts("a,b")
	want := []string{"a", ",", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestTokenTexts_floatFractionAfterDashCompound(t *testing.T) {
	t.Parallel()
	// Regression: float join must not break dash Rule2112 on the same line.
	got := tokenize.TokenTexts("что-то 1,5")
	want := []string{"что-то", "1,5"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
