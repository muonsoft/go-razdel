package sentenize

// defaultJoinRules matches upstream sentenize.RULES (trivial + sokr + list_item +
// close_quote + close_bracket + dash_right).
var defaultJoinRules = buildDefaultJoinRules()

func buildDefaultJoinRules() []func(*SentSplit) bool {
	n := len(TrivialJoinRules) + 3 + 1 + 2 + 1
	out := make([]func(*SentSplit) bool, 0, n)
	out = append(out, TrivialJoinRules...)
	out = append(out, RuleSokrLeft, RuleInsidePairSokr, RuleInitialsLeft)
	out = append(out, RuleListItem, RuleCloseQuote, RuleCloseBracket, RuleDashRight)
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
