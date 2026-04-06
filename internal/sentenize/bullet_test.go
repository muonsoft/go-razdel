package sentenize

import (
	"reflect"
	"testing"
)

func TestIsBulletToken_upstream(t *testing.T) {
	t.Parallel()
	cases := []struct {
		token string
		want  bool
	}{
		{"4", true},
		{".", true},
		{")", true},
		{"8", true},
		{"§", true},
		{"а", true},
		{"IV", true},
		{"XML", true}, // ROMAN charset IVXML
		{"iv", false},
		{"Hello", false},
		{"z", false},
	}
	for _, tc := range cases {
		if g := isBulletToken(tc.token); g != tc.want {
			t.Errorf("isBulletToken(%q)=%v want %v", tc.token, g, tc.want)
		}
	}
}

func TestRuleListItem(t *testing.T) {
	t.Parallel()
	buf := "4"
	sp := &SentSplit{Delimiter: ".", Buffer: &buf, Right: " Я"}
	if !RuleListItem(sp) {
		t.Fatal("4. + continuation should join")
	}
	buf8 := "8.1"
	sp2 := &SentSplit{Delimiter: ".", Buffer: &buf8, Right: " Зачем"}
	if !RuleListItem(sp2) {
		t.Fatal("8.1. bullet chain should join")
	}
	long := "123456789012345678901" // 21 runes
	sp3 := &SentSplit{Delimiter: ".", Buffer: &long, Right: " x"}
	if RuleListItem(sp3) {
		t.Fatal("buffer longer than BULLET_SIZE must not list_item join")
	}
	buf2 := "2"
	sp4 := &SentSplit{Delimiter: ")", Buffer: &buf2, Right: " отчуждать"}
	if !RuleListItem(sp4) {
		t.Fatal("2) list marker should join")
	}
	sp5 := &SentSplit{Delimiter: ";", Buffer: &buf, Right: " x"}
	if RuleListItem(sp5) {
		t.Fatal("semicolon delimiter must not trigger list_item")
	}
}

func TestRuleDashRight(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Delimiter: ".", Right: " - y"}
	if !RuleDashRight(sp) {
		t.Fatal("ASCII hyphen + lowercase should join")
	}
	sp2 := &SentSplit{Delimiter: ".", Right: " — тихо."}
	if !RuleDashRight(sp2) {
		t.Fatal("em dash + lowercase Cyrillic should join")
	}
	sp3 := &SentSplit{Delimiter: ".", Right: " — Тихо."}
	if RuleDashRight(sp3) {
		t.Fatal("em dash + uppercase must not join")
	}
	sp4 := &SentSplit{Delimiter: ".", Right: " - Y"}
	if RuleDashRight(sp4) {
		t.Fatal("hyphen + uppercase must not join")
	}
}

func TestSegment_upstreamUnitBulletsAndDialogue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		text string
		want []string
	}{
		{
			`4. Я присутствовал во время встречи`,
			[]string{`4. Я присутствовал во время встречи`},
		},
		{
			`IV. Гестационный сахарный диабет`,
			[]string{`IV. Гестационный сахарный диабет`},
		},
		{
			`§2. Нахождение оптимального объекта.`,
			[]string{`§2. Нахождение оптимального объекта.`},
		},
		{
			`8.1. Зачем нужны эти классы?`,
			[]string{`8.1. Зачем нужны эти классы?`},
		},
		{
			`в данной квартире; 2) отчуждать свою долю`,
			[]string{`в данной квартире;`, `2) отчуждать свою долю`},
		},
		{
			`- "Так в чем же дело?" - "Не ра-ду-ют".`,
			[]string{`- "Так в чем же дело?"`, `- "Не ра-ду-ют".`},
		},
		{
			`— Ты ей скажи, что я ей гостинца дам. — А мне дашь?`,
			[]string{`— Ты ей скажи, что я ей гостинца дам.`, `— А мне дашь?`},
		},
		{
			`x. - y`,
			[]string{`x. - y`},
		},
		{
			`x. - Y`,
			[]string{`x.`, `- Y`},
		},
	}
	for _, tc := range cases {
		parts := SentSplitterParts(tc.text, DefaultWindow)
		got := PostStrip(Segment(parts, JoinDefault))
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("text=%q\ngot  %#v\nwant %#v", tc.text, got, tc.want)
		}
	}
}

func TestSegment_regressionBoundAndSokrStillSplit(t *testing.T) {
	t.Parallel()
	cases := []struct {
		text string
		want []string
	}{
		{
			`словам, "не будет точно". "Возможно, у нас`,
			[]string{`словам, "не будет точно".`, `"Возможно, у нас`},
		},
		{
			`И т. д. и т. п. В общем, вся газета`,
			[]string{`И т. д. и т. п.`, `В общем, вся газета`},
		},
	}
	for _, tc := range cases {
		parts := SentSplitterParts(tc.text, DefaultWindow)
		got := PostStrip(Segment(parts, JoinDefault))
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("text=%q\ngot  %#v\nwant %#v", tc.text, got, tc.want)
		}
	}
}
