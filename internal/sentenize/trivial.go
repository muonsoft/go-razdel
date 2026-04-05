package sentenize

import (
	"regexp"
	"strings"
	"unicode"
)

// smilePrefixRE matches upstream SMILE_PREFIX: optional leading whitespace then SMILES.
var smilePrefixRE = regexp.MustCompile(`(?s)^\s*[=:;]-?[)(]{1,3}`)

// TrivialJoinRules is upstream RULES[:4] order: empty_side, no_space_prefix, lower_right, delimiter_right.
var TrivialJoinRules = []func(*SentSplit) bool{
	RuleEmptySide,
	RuleNoSpacePrefix,
	RuleLowerRight,
	RuleDelimiterRight,
}

// JoinTrivial returns true when upstream would JOIN for the trivial rule layer only.
func JoinTrivial(s *SentSplit) bool {
	for _, r := range TrivialJoinRules {
		if r(s) {
			return true
		}
	}
	return false
}

// RuleEmptySide mirrors upstream empty_side.
func RuleEmptySide(s *SentSplit) bool {
	return s.LeftToken() == "" || s.RightToken() == ""
}

// RuleNoSpacePrefix mirrors upstream no_space_prefix.
func RuleNoSpacePrefix(s *SentSplit) bool {
	return !s.RightSpacePrefix()
}

// RuleLowerRight mirrors upstream lower_right (is_lower_alpha on right_token).
func RuleLowerRight(s *SentSplit) bool {
	return isLowerAlphaToken(s.RightToken())
}

// RuleDelimiterRight mirrors upstream delimiter_right.
func RuleDelimiterRight(s *SentSplit) bool {
	right := s.RightToken()
	if strings.Contains(genericQuotes, right) {
		return false
	}
	if strings.Contains(delimiters, right) {
		return true
	}
	return smilePrefixRE.MatchString(s.Right)
}

func isLowerAlphaToken(token string) bool {
	if token == "" {
		return false
	}
	for _, r := range token {
		if !unicode.IsLetter(r) {
			return false
		}
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}
