package sentenize

import (
	"testing"
)

func TestRuleEmptySide(t *testing.T) {
	if !RuleEmptySide(&SentSplit{Left: "", Delimiter: ".", Right: "x"}) {
		t.Fatal("empty left_token should join")
	}
	if !RuleEmptySide(&SentSplit{Left: "a", Delimiter: ".", Right: ""}) {
		t.Fatal("empty right_token should join")
	}
	if RuleEmptySide(&SentSplit{Left: "a", Delimiter: ".", Right: "b"}) {
		t.Fatal("both sides should not trigger empty_side")
	}
}

func TestRuleNoSpacePrefix(t *testing.T) {
	if !RuleNoSpacePrefix(&SentSplit{Left: "a", Delimiter: ".", Right: "b"}) {
		t.Fatal("no leading space on right")
	}
	if RuleNoSpacePrefix(&SentSplit{Left: "a", Delimiter: ".", Right: " b"}) {
		t.Fatal("space prefix should not trigger no_space_prefix")
	}
}

func TestRuleLowerRight(t *testing.T) {
	if !RuleLowerRight(&SentSplit{Left: "a", Delimiter: ".", Right: "b"}) {
		t.Fatal("lowercase letter should join")
	}
	if !RuleLowerRight(&SentSplit{Left: "a", Delimiter: ".", Right: "привет"}) {
		t.Fatal("cyrillic lowercase should join")
	}
	if RuleLowerRight(&SentSplit{Left: "a", Delimiter: ".", Right: "B"}) {
		t.Fatal("uppercase should not join")
	}
	if RuleLowerRight(&SentSplit{Left: "a", Delimiter: ".", Right: "9"}) {
		t.Fatal("digit should not join")
	}
}

func TestRuleDelimiterRight(t *testing.T) {
	if RuleDelimiterRight(&SentSplit{Left: "a", Delimiter: ".", Right: ` "x"`}) {
		t.Fatal("generic quote token should not join via delimiter_right")
	}
	if !RuleDelimiterRight(&SentSplit{Left: "a", Delimiter: ".", Right: ";x"}) {
		t.Fatal("semicolon follower should join")
	}
	if !RuleDelimiterRight(&SentSplit{Left: "a", Delimiter: ".", Right: ":-)x"}) {
		t.Fatal("smile prefix should join")
	}
	if !RuleDelimiterRight(&SentSplit{Left: "a", Delimiter: ".", Right: " :-)x"}) {
		t.Fatal("smile after spaces should join")
	}
}

func TestTrivialRuleOrderMatchesUpstream(t *testing.T) {
	// If order were lower_right before no_space_prefix, "a. B" could be split differently
	// depending on implementation; upstream runs no_space_prefix second.
	sp := &SentSplit{Left: "Hello", Delimiter: ".", Right: " World"}
	if JoinTrivial(sp) {
		t.Fatal("Hello. World should not trivial-join (space + uppercase W)")
	}
}

func TestJoinTrivial_firstRuleWins(t *testing.T) {
	// empty_side must run before no_space_prefix: both could apply conceptually on edge cases;
	// empty right is caught by empty_side.
	sp := &SentSplit{Left: "привет", Delimiter: ".", Right: ""}
	if !JoinTrivial(sp) {
		t.Fatal("expected join for trailing delimiter / empty right")
	}
}
