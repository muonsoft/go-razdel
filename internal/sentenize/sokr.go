package sentenize

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	sokrs         map[string]struct{}
	headSokrs     map[string]struct{}
	pairSokrs     map[[2]string]struct{}
	headPairSokrs map[[2]string]struct{}
)

// initialsSet mirrors upstream INITIALS in sokr.py (lowercase keys).
var initialsSet = map[string]struct{}{
	"дж": {},
	"ed": {},
	"вс": {},
}

func init() {
	sokrs = wordSet(sokrWordsRaw)
	headSokrs = wordSet(headSokrWordsRaw)
	pairSokrs = pairSet(pairSokrPairs)
	headPairSokrs = pairSet(headPairSokrPairs)
}

func wordSet(raw string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, w := range strings.Fields(raw) {
		m[w] = struct{}{}
	}
	return m
}

func pairSet(pairs [][2]string) map[[2]string]struct{} {
	m := make(map[[2]string]struct{}, len(pairs))
	for _, p := range pairs {
		m[p] = struct{}{}
	}
	return m
}

// isSokr mirrors upstream is_sokr (sentenize.py).
func isSokr(token string) bool {
	if token == "" {
		return true
	}
	if isAllDecimalDigits(token) {
		return true
	}
	if !isAllLetters(token) {
		return true
	}
	return isAllLowerLetters(token)
}

func isAllDecimalDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isAllLetters(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isAllLowerLetters(s string) bool {
	hasCased := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			if !unicode.IsLower(r) {
				return false
			}
			hasCased = true
		}
	}
	return hasCased
}

// RuleSokrLeft mirrors upstream sokr_left.
func RuleSokrLeft(s *SentSplit) bool {
	if s.Delimiter != "." {
		return false
	}
	right := s.RightToken()
	if a, b, ok := s.LeftPairSokr(); ok {
		lp := [2]string{strings.ToLower(a), strings.ToLower(b)}
		if _, h := headPairSokrs[lp]; h {
			return true
		}
		if _, p := pairSokrs[lp]; p {
			return isSokr(right)
		}
	}
	left := strings.ToLower(s.LeftToken())
	if _, h := headSokrs[left]; h {
		return true
	}
	if _, m := sokrs[left]; m {
		return isSokr(right)
	}
	return false
}

// RuleInsidePairSokr mirrors upstream inside_pair_sokr.
func RuleInsidePairSokr(s *SentSplit) bool {
	if s.Delimiter != "." {
		return false
	}
	lp := [2]string{strings.ToLower(s.LeftToken()), strings.ToLower(s.RightToken())}
	_, ok := pairSokrs[lp]
	return ok
}

// RuleInitialsLeft mirrors upstream initials_left.
func RuleInitialsLeft(s *SentSplit) bool {
	if s.Delimiter != "." {
		return false
	}
	left := s.LeftToken()
	r, sz := utf8.DecodeRuneInString(left)
	if sz == len(left) && unicode.IsLetter(r) && unicode.IsUpper(r) {
		return true
	}
	_, ok := initialsSet[strings.ToLower(left)]
	return ok
}
