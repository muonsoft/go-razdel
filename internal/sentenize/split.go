// Package sentenize implements sentence split primitives aligned with upstream
// third_party/razdel/razdel/segmenters/sentenize.py (SentSplit / SentSplitter).
package sentenize

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

// DefaultWindow is the upstream SentSplitter default (10 Unicode code points).
const DefaultWindow = 10

// delimiterRE matches upstream DELIMITER: smileys or punctuation/closing delimiters.
// See sentenize.DELIMITER in third_party/razdel.
var delimiterRE = regexp.MustCompile(`([=:;]-?[)(]{1,3}|[.?!…;"\x{201e}'\x{bb}\x{201d}\x{2019}\)\]\}])`)

var (
	firstTokenRE = regexp.MustCompile(`(?s)^\s*([\p{L}\p{M}\p{Pc}\p{Nl}\p{No}]+|\p{Nd}+|[^\p{L}\p{N}\p{M}\p{Pc}\p{Z}])`)
	lastTokenRE  = regexp.MustCompile(`(?s)([\p{L}\p{M}\p{Pc}\p{Nl}\p{No}]+|\p{Nd}+|[^\p{L}\p{N}\p{M}\p{Pc}\p{Z}])\s*$`)
	wordRE       = regexp.MustCompile(`(?s)([\p{L}\p{M}\p{Pc}\p{Nl}\p{No}]+|\p{Nd}+)`)
	tokenRE      = regexp.MustCompile(`[\p{L}\p{M}\p{Pc}\p{Nl}\p{No}]+|\p{Nd}+|[^\p{L}\p{N}\p{M}\p{Pc}\p{Z}]`)
	pairSokrRE   = regexp.MustCompile(`(?s)([\p{L}\p{N}\p{M}\p{Pc}])\s*\.\s*([\p{L}\p{N}\p{M}\p{Pc}])\s*$`)
	intSokrRE    = regexp.MustCompile(`(?s)\p{Nd}+\s*-?\s*([\p{L}\p{N}\p{M}\p{Pc}]+)\s*$`)
)

// SentSplit mirrors upstream SentSplit(left, delimiter, right, buffer=None).
type SentSplit struct {
	Left, Delimiter, Right string
	Buffer                 *string

	leftToken       string
	leftTokenSource string
	leftTokenCached bool

	rightToken       string
	rightTokenSource string
	rightTokenCached bool

	rightWord       string
	rightWordSource string
	rightWordCached bool

	bufferTokens       []string
	bufferTokensSource *string
	bufferTokensCached bool

	bufferFirstToken       string
	bufferFirstTokenSource *string
	bufferFirstTokenCached bool
}

// RightSpacePrefix mirrors SentSplit.right_space_prefix (upstream SPACE_PREFIX).
func (s *SentSplit) RightSpacePrefix() bool {
	if s.Right == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s.Right)
	return unicode.IsSpace(r)
}

// LeftSpaceSuffix mirrors SentSplit.left_space_suffix (upstream SPACE_SUFFIX).
func (s *SentSplit) LeftSpaceSuffix() bool {
	if s.Left == "" {
		return false
	}
	r, _ := utf8.DecodeLastRuneInString(s.Left)
	return unicode.IsSpace(r)
}

// RightToken mirrors SentSplit.right_token (upstream FIRST_TOKEN).
func (s *SentSplit) RightToken() string {
	if s.rightTokenCached && s.rightTokenSource == s.Right {
		return s.rightToken
	}
	s.rightTokenSource = s.Right
	s.rightTokenCached = true
	m := firstTokenRE.FindStringSubmatch(s.Right)
	if m == nil {
		s.rightToken = ""
		return ""
	}
	s.rightToken = m[1]
	return s.rightToken
}

// LeftToken mirrors SentSplit.left_token (upstream LAST_TOKEN).
func (s *SentSplit) LeftToken() string {
	if s.leftTokenCached && s.leftTokenSource == s.Left {
		return s.leftToken
	}
	s.leftTokenSource = s.Left
	s.leftTokenCached = true
	m := lastTokenRE.FindStringSubmatch(s.Left)
	if m == nil {
		s.leftToken = ""
		return ""
	}
	s.leftToken = m[1]
	return s.leftToken
}

// LeftPairSokr mirrors SentSplit.left_pair_sokr (upstream PAIR_SOKR).
func (s *SentSplit) LeftPairSokr() (a, b string, ok bool) {
	loc := pairSokrRE.FindStringSubmatchIndex(s.Left)
	if loc == nil {
		return "", "", false
	}
	return s.Left[loc[2]:loc[3]], s.Left[loc[4]:loc[5]], true
}

// LeftIntSokr mirrors SentSplit.left_int_sokr (upstream INT_SOKR).
func (s *SentSplit) LeftIntSokr() (word string, ok bool) {
	loc := intSokrRE.FindStringSubmatchIndex(s.Left)
	if loc == nil {
		return "", false
	}
	return s.Left[loc[2]:loc[3]], true
}

// RightWord mirrors SentSplit.right_word (upstream WORD).
func (s *SentSplit) RightWord() string {
	if s.rightWordCached && s.rightWordSource == s.Right {
		return s.rightWord
	}
	s.rightWordSource = s.Right
	s.rightWordCached = true
	m := wordRE.FindStringSubmatch(s.Right)
	if m == nil {
		s.rightWord = ""
		return ""
	}
	s.rightWord = m[1]
	return s.rightWord
}

// BufferTokens mirrors SentSplit.buffer_tokens (upstream TOKEN.findall); unset buffer returns nil.
func (s *SentSplit) BufferTokens() []string {
	if s.Buffer == nil {
		return nil
	}
	if s.bufferTokensCached && s.bufferTokensSource == s.Buffer {
		return s.bufferTokens
	}
	s.bufferTokensSource = s.Buffer
	s.bufferTokensCached = true
	s.bufferTokens = tokenRE.FindAllString(*s.Buffer, -1)
	return s.bufferTokens
}

// BufferFirstToken mirrors SentSplit.buffer_first_token; unset buffer returns "".
func (s *SentSplit) BufferFirstToken() string {
	if s.Buffer == nil {
		return ""
	}
	if s.bufferFirstTokenCached && s.bufferFirstTokenSource == s.Buffer {
		return s.bufferFirstToken
	}
	s.bufferFirstTokenSource = s.Buffer
	s.bufferFirstTokenCached = true
	m := firstTokenRE.FindStringSubmatch(*s.Buffer)
	if m == nil {
		s.bufferFirstToken = ""
		return ""
	}
	s.bufferFirstToken = m[1]
	return s.bufferFirstToken
}

func onlyWhitespace(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// runePrefixBefore returns the last maxRunes runes of s[:endByte].
func runePrefixBefore(s string, endByte, maxRunes int) string {
	if endByte <= 0 || maxRunes <= 0 {
		return ""
	}
	if endByte > len(s) {
		endByte = len(s)
	}
	start := endByte
	for range maxRunes {
		if start <= 0 {
			break
		}
		_, width := utf8.DecodeLastRuneInString(s[:start])
		if width == 0 {
			break
		}
		start -= width
	}
	return s[start:endByte]
}

// runeSuffixAfter returns the first maxRunes runes of s[startByte:].
func runeSuffixAfter(s string, startByte, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	if startByte < 0 {
		startByte = 0
	}
	if startByte >= len(s) {
		return ""
	}
	stop := startByte
	for range maxRunes {
		if stop >= len(s) {
			break
		}
		_, width := utf8.DecodeRuneInString(s[stop:])
		if width == 0 {
			break
		}
		stop += width
	}
	return s[startByte:stop]
}

// SentPart is one element of the interleaved stream from SentSplitter.__call__:
// either a text chunk (Split == nil) or a split point (Split != nil).
type SentPart struct {
	Text  string
	Split *SentSplit
}

// SentSplitterParts returns the interleaved stream from upstream SentSplitter.__call__:
// text chunks and *SentSplit events. Window is in Unicode code points (Python str slice).
// For whitespace-only input, returns nil (upstream yields nothing).
func SentSplitterParts(text string, window int) []SentPart {
	if onlyWhitespace(text) {
		return nil
	}
	if window <= 0 {
		window = DefaultWindow
	}
	idx := delimiterRE.FindAllStringIndex(text, -1)
	if len(idx) == 0 {
		return []SentPart{{Text: text}}
	}
	out := make([]SentPart, 0, len(idx)*2+1)
	prev := 0
	for _, loc := range idx {
		start, stop := loc[0], loc[1]
		delim := text[start:stop]
		out = append(out, SentPart{Text: text[prev:start]})
		out = append(out, SentPart{Split: &SentSplit{
			Left:      runePrefixBefore(text, start, window),
			Delimiter: delim,
			Right:     runeSuffixAfter(text, stop, window),
		}})
		prev = stop
	}
	out = append(out, SentPart{Text: text[prev:]})
	return out
}
