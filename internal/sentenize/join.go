package sentenize

// defaultJoinRules matches upstream sentenize.RULES through close_bracket, excluding
// list_item (before close_quote) and dash_right (after close_bracket) — see T015.
var defaultJoinRules = buildDefaultJoinRules()

func buildDefaultJoinRules() []func(*SentSplit) bool {
	n := len(TrivialJoinRules) + 3 + 2
	out := make([]func(*SentSplit) bool, 0, n)
	out = append(out, TrivialJoinRules...)
	out = append(out, RuleSokrLeft, RuleInsidePairSokr, RuleInitialsLeft)
	// T015: insert RuleListItem here (before close_quote), append RuleDashRight after RuleCloseBracket.
	out = append(out, RuleCloseQuote, RuleCloseBracket)
	return out
}

// JoinDefault applies defaultJoinRules in upstream order; first matching rule wins.
func JoinDefault(s *SentSplit) bool {
	for _, r := range defaultJoinRules {
		if r(s) {
			return true
		}
	}
	return false
}
