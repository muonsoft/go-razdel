package sentenize

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Upstream sentenize: ROMAN, BULLET_CHARS, BULLET_BOUNDS, BULLET_SIZE.
const bulletSize = 20

var romanRE = regexp.MustCompile(`^[IVXML]+$`)

var bulletCharSet = stringToRuneSet("§абвгдеabcdef")

func stringToRuneSet(s string) map[rune]struct{} {
	m := make(map[rune]struct{}, utf8.RuneCountInString(s))
	for _, r := range s {
		m[r] = struct{}{}
	}
	return m
}

// isBulletToken mirrors upstream is_bullet (sentenize.py).
func isBulletToken(token string) bool {
	if token == "" {
		return false
	}
	if stringIsDigit(token) {
		return true
	}
	if strings.Contains(".)", token) {
		return true
	}
	if tok := strings.ToLower(token); utf8.RuneCountInString(tok) == 1 {
		r, _ := utf8.DecodeRuneInString(tok)
		if _, ok := bulletCharSet[r]; ok {
			return true
		}
	}
	return romanRE.MatchString(token)
}

func stringIsDigit(s string) bool {
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

// RuleListItem mirrors upstream list_item.
func RuleListItem(s *SentSplit) bool {
	if !strings.Contains(".)", s.Delimiter) {
		return false
	}
	if s.Buffer == nil {
		return false
	}
	if utf8.RuneCountInString(*s.Buffer) > bulletSize {
		return false
	}
	for _, tok := range s.BufferTokens() {
		if !isBulletToken(tok) {
			return false
		}
	}
	return true
}

// RuleDashRight mirrors upstream dash_right.
func RuleDashRight(s *SentSplit) bool {
	rt := s.RightToken()
	if rt == "" || !strings.Contains(Dashes, rt) {
		return false
	}
	rw := s.RightWord()
	if rw == "" {
		return false
	}
	return isLowerAlphaToken(rw)
}
