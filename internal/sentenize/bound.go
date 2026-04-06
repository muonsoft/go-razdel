package sentenize

import "strings"

// leftTokenInEndings mirrors upstream `left in ENDINGS` when ENDINGS is a str:
// substring membership, with Python's special case that "" is contained in any str.
func leftTokenInEndings(s *SentSplit) bool {
	left := s.LeftToken()
	if left == "" {
		return true
	}
	return strings.Contains(Endings, left)
}

// closeBound mirrors upstream close_bound: JOIN unless left_token is in ENDINGS (Python str containment).
func closeBound(s *SentSplit) bool {
	return !leftTokenInEndings(s)
}

// RuleCloseQuote mirrors upstream close_quote (delimiter in QUOTES / CLOSE_QUOTES / GENERIC_QUOTES).
func RuleCloseQuote(s *SentSplit) bool {
	d := s.Delimiter
	if !strings.Contains(Quotes, d) {
		return false
	}
	if strings.Contains(CloseQuotes, d) {
		return closeBound(s)
	}
	if strings.Contains(GenericQuotes, d) {
		if !s.LeftSpaceSuffix() {
			return closeBound(s)
		}
		return true
	}
	return false
}

// RuleCloseBracket mirrors upstream close_bracket.
func RuleCloseBracket(s *SentSplit) bool {
	if !strings.Contains(CloseBrackets, s.Delimiter) {
		return false
	}
	return closeBound(s)
}
