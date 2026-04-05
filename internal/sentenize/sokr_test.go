package sentenize

import (
	"testing"
)

func TestSokrDictionaries_integrity(t *testing.T) {
	t.Parallel()
	if len(sokrs) != 137 {
		t.Fatalf("SOKRS: want 137 keys (upstream), got %d", len(sokrs))
	}
	if len(headSokrs) != 103 {
		t.Fatalf("HEAD_SOKRS: want 103 keys (upstream), got %d", len(headSokrs))
	}
	if len(pairSokrs) != 24 {
		t.Fatalf("PAIR_SOKRS 2-tuples: want 24 (upstream), got %d", len(pairSokrs))
	}
	if len(headPairSokrs) != 9 {
		t.Fatalf("HEAD_PAIR_SOKRS: want 9 keys (upstream), got %d", len(headPairSokrs))
	}

	keyCases := []struct {
		name string
		m    map[string]struct{}
		key  string
	}{
		{"SOKRS т", sokrs, "т"},
		{"SOKRS проц", sokrs, "проц"},
		{"HEAD оз", headSokrs, "оз"},
		{"initials дж", initialsSet, "дж"},
	}
	for _, tc := range keyCases {
		if _, ok := tc.m[tc.key]; !ok {
			t.Fatalf("%s: missing key %q", tc.name, tc.key)
		}
	}
	pair := [2]string{"т", "д"}
	if _, ok := pairSokrs[pair]; !ok {
		t.Fatal("PAIR_SOKRS missing (т, д)")
	}
	if _, ok := headPairSokrs[[2]string{"т", "е"}]; !ok {
		t.Fatal("HEAD_PAIR_SOKRS missing (т, е)")
	}
}

func TestRuleInsidePairSokr(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "и т", Delimiter: ".", Right: " д. далее"}
	if !RuleInsidePairSokr(sp) {
		t.Fatal("т. д. continuation should join")
	}
	sp2 := &SentSplit{Left: "слово т", Delimiter: ".", Right: " Далее"}
	if RuleInsidePairSokr(sp2) {
		t.Fatal("(т, Д) lower mismatch should not match PAIR_SOKRS alone — wrong rule")
	}
}

func TestRuleSokrLeft_headPairUnconditional(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "к.п", Delimiter: ".", Right: " Н"}
	if !RuleSokrLeft(sp) {
		t.Fatal("HEAD_PAIR (к,п) should join even when right is uppercase")
	}
}

func TestRuleSokrLeft_pairRequiresSokrRight(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "т.д", Delimiter: ".", Right: " Далее"}
	if RuleSokrLeft(sp) {
		t.Fatal("TAIL_PAIR (т,д) must not join before uppercase continuation")
	}
	sp2 := &SentSplit{Left: "т.д", Delimiter: ".", Right: " далее"}
	if !RuleSokrLeft(sp2) {
		t.Fatal("TAIL_PAIR (т,д) should join before lowercase continuation")
	}
}

func TestRuleSokrLeft_headAbbrev(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "оз", Delimiter: ".", Right: " Селяха"}
	if !RuleSokrLeft(sp) {
		t.Fatal("оз. as HEAD_SOKR should join")
	}
}

func TestRuleInitialsLeft(t *testing.T) {
	t.Parallel()
	if !RuleInitialsLeft(&SentSplit{Left: "Л", Delimiter: ".", Right: " В"}) {
		t.Fatal("single-letter uppercase should join")
	}
	if !RuleInitialsLeft(&SentSplit{Left: "Дж", Delimiter: ".", Right: " Ф"}) {
		t.Fatal("Дж in INITIALS should join")
	}
	if RuleInitialsLeft(&SentSplit{Left: "аб", Delimiter: ".", Right: " в"}) {
		t.Fatal("non-initial two-letter lowercase should not match initials_left")
	}
}

func TestJoinDefault_sokrAfterTrivial(t *testing.T) {
	t.Parallel()
	// delimiter_right would join if right starts with ";"; sokr should not be needed.
	sp := &SentSplit{Left: "т.д", Delimiter: ".", Right: " далее"}
	if !JoinDefault(sp) {
		t.Fatal("expected join via sokr_left after trivial rules pass")
	}
}

func TestJoinDefault_regression_plainWord(t *testing.T) {
	t.Parallel()
	sp := &SentSplit{Left: "привет", Delimiter: ".", Right: " Мир"}
	if JoinDefault(sp) {
		t.Fatal("plain word before dot should not join into next sentence")
	}
}
