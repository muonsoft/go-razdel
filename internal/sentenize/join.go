package sentenize

// defaultJoinRules matches upstream sentenize.RULES through initials_left
// (empty_side … initials_left), excluding list_item, quote/bracket bounds, dash_right.
var defaultJoinRules = buildDefaultJoinRules()

func buildDefaultJoinRules() []func(*SentSplit) bool {
	n := len(TrivialJoinRules) + 3
	out := make([]func(*SentSplit) bool, 0, n)
	out = append(out, TrivialJoinRules...)
	out = append(out, RuleSokrLeft, RuleInsidePairSokr, RuleInitialsLeft)
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
