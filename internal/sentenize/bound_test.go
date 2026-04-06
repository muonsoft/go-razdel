package sentenize

import (
	"reflect"
	"testing"
)

func TestRuleCloseQuote_closeQuoteUsesCloseBound(t *testing.T) {
	t.Parallel()
	// Closing » after non-ending token → join (merge across »).
	sp := &SentSplit{Left: "слово", Delimiter: "\u00bb", Right: " далее"}
	if !RuleCloseQuote(sp) {
		t.Fatal("expected join when left_token is not an ending")
	}
	// Ellipsis before » — left_token … is in ENDINGS → do not join from this rule.
	sp2 := &SentSplit{Left: "кто они такие\u2026 ", Delimiter: "\u00bb", Right: ""}
	if RuleCloseQuote(sp2) {
		t.Fatal("expected no join when left_token is ending before »")
	}
}

func TestRuleCloseQuote_genericQuoteSpaceSuffixAlwaysJoins(t *testing.T) {
	t.Parallel()
	// GENERIC_QUOTES + left_space_suffix → JOIN (upstream), regardless of close_bound.
	sp := &SentSplit{Left: "текст ", Delimiter: `"`, Right: "продолжение"}
	if !RuleCloseQuote(sp) {
		t.Fatal("generic quote with left_space_suffix should join")
	}
}

func TestRuleCloseQuote_genericQuoteNoSpaceUsesCloseBound(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: `слово"не`, Delimiter: `"`, Right: " далее"}
	if !RuleCloseQuote(sp) {
		t.Fatal(`expected join: left_token before " is not an ending`)
	}
	sp2 := &SentSplit{Left: `текст.`, Delimiter: `"`, Right: ` Далее`}
	if RuleCloseQuote(sp2) {
		t.Fatal(`expected no join when generic quote follows ending: left_token is .`)
	}
}

func TestRuleCloseBracket(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "слово", Delimiter: ")", Right: " далее"}
	if !RuleCloseBracket(sp) {
		t.Fatal(") after word should join")
	}
	sp2 := &SentSplit{Left: "текст.", Delimiter: ")", Right: " далее"}
	if RuleCloseBracket(sp2) {
		t.Fatal(") after ending should not join from close_bracket")
	}
}

func TestRuleCloseQuote_openQuoteDoesNotMatchBranches(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "x", Delimiter: "\u00ab", Right: " y"}
	if RuleCloseQuote(sp) {
		t.Fatal("open quote delimiter should not trigger close_quote join/continue mismatch")
	}
}

func TestSegment_boundParityUpstream(t *testing.T) {
	// Expected chunks match third_party/razdel/razdel/tests/test_sentenize.py BOUND (UNIT) cases
	// under the same rule set omitting list_item and dash_right (T014 scope).
	cases := []struct {
		text string
		want []string
	}{
		{
			`словам, "не будет точно". "Возможно, у нас`,
			[]string{`словам, "не будет точно".`, `"Возможно, у нас`},
		},
		{
			`Брось!.." Связываться не хотелось`,
			[]string{`Брось!.."`, `Связываться не хотелось`},
		},
		{
			`Peter Goldreich,Scott Tremaine (1979). «Относительно теории колец Урана».`,
			[]string{`Peter Goldreich,Scott Tremaine (1979).`, `«Относительно теории колец Урана».`},
		},
		{
			`Это чудовищные риски. "Яндекс" попал под удар`,
			[]string{`Это чудовищные риски.`, `"Яндекс" попал под удар`},
		},
		{
			`кто они такие… »`,
			[]string{`кто они такие… »`},
		},
	}
	for _, tc := range cases {
		parts := SentSplitterParts(tc.text, DefaultWindow)
		raw := Segment(parts, JoinDefault)
		got := PostStrip(raw)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("text=%q\ngot  %#v\nwant %#v", tc.text, got, tc.want)
		}
	}
}
