package tokenize

import (
	"regexp"
	"strings"
)

// join mirrors upstream razdel.rule.JOIN (non-zero action).
const join = 1

var smileRE = regexp.MustCompile(`^[=:;]-?[)(]{1,3}$`)

// splitSpace mirrors upstream tokenize.split_space: non-empty delimiter always
// yields SPLIT (tokens do not join across spaces or other gaps between atoms).
// See third_party/razdel/razdel/segmenters/tokenize.py.
func splitSpace(s *tokenSplit) bool {
	return s.Delimiter != ""
}

// shouldJoin mirrors TokenSegmenter.join when only join rules apply after trivial
// split: join iff delimiter is empty and a rule returns JOIN.
func shouldJoin(s *tokenSplit) bool {
	if splitSpace(s) {
		return false
	}
	return joinSplit(s)
}

// tokenSplit carries atom windows for join rules (upstream TokenSplit).
type tokenSplit struct {
	leftAtoms  []Atom
	rightAtoms []Atom
	Delimiter  string // bytes between previous atom end and current atom start
	buffer     string // filled during segment pass
}

func (s *tokenSplit) left1() *Atom  { return &s.leftAtoms[len(s.leftAtoms)-1] }
func (s *tokenSplit) left2() *Atom  { return &s.leftAtoms[len(s.leftAtoms)-2] }
func (s *tokenSplit) right1() *Atom { return &s.rightAtoms[0] }
func (s *tokenSplit) right2() *Atom { return &s.rightAtoms[1] }

// Left and Right mirror Split.left / Split.right (neighbor atom texts).
func (s *tokenSplit) Left() string  { return s.left1().Text }
func (s *tokenSplit) Right() string { return s.right1().Text }

func dashDelimiterPiece(piece string) bool {
	if piece == "" {
		return false
	}
	for _, r := range piece {
		if !strings.ContainsRune(Dashes, r) {
			return false
		}
	}
	return true
}

func underscoreDelimiterPiece(piece string) bool {
	return piece == "_"
}

func floatDelimiterPiece(piece string) bool {
	if len(piece) != 1 {
		return false
	}
	switch piece[0] {
	case '.', ',':
		return true
	default:
		return false
	}
}

func fractionDelimiterPiece(piece string) bool {
	if len(piece) != 1 {
		return false
	}
	switch piece[0] {
	case '/', '\\':
		return true
	default:
		return false
	}
}

// rule2112Action returns join or split or 0 if rule does not apply.
func dashRule2112(s *tokenSplit) int {
	var left, right *Atom
	switch {
	case dashDelimiterPiece(s.Left()):
		if len(s.leftAtoms) < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case dashDelimiterPiece(s.Right()):
		if len(s.rightAtoms) < 2 {
			return 0
		}
		left, right = s.left1(), s.right2()
	default:
		return 0
	}
	if left.Type == PUNCT || right.Type == PUNCT {
		return 0
	}
	return join
}

func underscoreRule2112(s *tokenSplit) int {
	var left, right *Atom
	switch {
	case underscoreDelimiterPiece(s.Left()):
		if len(s.leftAtoms) < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case underscoreDelimiterPiece(s.Right()):
		if len(s.rightAtoms) < 2 {
			return 0
		}
		left, right = s.left1(), s.right2()
	default:
		return 0
	}
	if left.Type == PUNCT || right.Type == PUNCT {
		return 0
	}
	return join
}

func floatRule2112(s *tokenSplit) int {
	var left, right *Atom
	switch {
	case floatDelimiterPiece(s.Left()):
		if len(s.leftAtoms) < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case floatDelimiterPiece(s.Right()):
		if len(s.rightAtoms) < 2 {
			return 0
		}
		left, right = s.left1(), s.right2()
	default:
		return 0
	}
	if left.Type == INT && right.Type == INT {
		return join
	}
	return 0
}

func fractionRule2112(s *tokenSplit) int {
	var left, right *Atom
	switch {
	case fractionDelimiterPiece(s.Left()):
		if len(s.leftAtoms) < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case fractionDelimiterPiece(s.Right()):
		if len(s.rightAtoms) < 2 {
			return 0
		}
		left, right = s.left1(), s.right2()
	default:
		return 0
	}
	if left.Type == INT && right.Type == INT {
		return join
	}
	return 0
}

func punctRule(s *tokenSplit) int {
	if s.left1().Type != PUNCT || s.right1().Type != PUNCT {
		return 0
	}
	left := s.Left()
	right := s.Right()
	combined := s.buffer + right
	if smileRE.MatchString(combined) {
		return join
	}
	// Lookahead: `:-)` / `;-)` need `:`/`;` + `-` joined before the closing `)`
	// matches SMILE as a whole (upstream #17; Python tokenize still splits here).
	if len(s.rightAtoms) >= 2 && smileRE.MatchString(combined+s.right2().Text) {
		return join
	}
	if strings.ContainsAny(Endings, left) && strings.ContainsAny(Endings, right) {
		return join
	}
	if left+right == "--" || left+right == "**" {
		return join
	}
	return 0
}

func otherRule(s *tokenSplit) int {
	l := s.left1().Type
	r := s.right1().Type
	if l == OTHER && (r == OTHER || r == RU || r == LAT) {
		if otherShouldSplitEmojiFromWord(*s.left1(), *s.right1()) {
			return 0
		}
		return join
	}
	if (l == OTHER || l == RU || l == LAT) && r == OTHER {
		if otherShouldSplitEmojiFromWord(*s.left1(), *s.right1()) {
			return 0
		}
		return join
	}
	return 0
}

func yahooRule(s *tokenSplit) int {
	if strings.EqualFold(s.left1().Text, "yahoo") && s.Right() == "!" {
		return join
	}
	return 0
}

func joinSplit(s *tokenSplit) bool {
	rules := []func(*tokenSplit) int{
		dashRule2112,
		underscoreRule2112,
		floatRule2112,
		fractionRule2112,
		punctRule,
		otherRule,
		yahooRule,
	}
	for _, rule := range rules {
		if a := rule(s); a != 0 {
			return a == join
		}
	}
	return false
}

// tokenPart is one element of the interleaved stream from TokenSplitter.__call__:
// either a text chunk (split == nil) or a split point (split != nil).
type tokenPart struct {
	text  string
	split *tokenSplit
}

// segmentStrings applies upstream TokenSegmenter.segment to the interleaved
// split/text stream from tokenSplitterParts.
func segmentStrings(parts []tokenPart) []string {
	if len(parts) == 0 {
		return nil
	}
	buffer := parts[0].text
	var out []string
	for i := 1; i < len(parts); i += 2 {
		sp := parts[i].split
		right := parts[i+1].text
		sp.buffer = buffer
		if shouldJoin(sp) {
			buffer += right
		} else {
			out = append(out, buffer)
			buffer = right
		}
	}
	out = append(out, buffer)
	return out
}

const window = 3

// tokenSplitterParts builds the interleaved text/split/text stream (upstream TokenSplitter.__call__).
func tokenSplitterParts(text string, atoms []Atom) []tokenPart {
	if len(atoms) == 0 {
		return nil
	}
	parts := make([]tokenPart, 0, len(atoms)*2)
	parts = append(parts, tokenPart{text: atoms[0].Text})
	for i := 1; i < len(atoms); i++ {
		prev := atoms[i-1]
		cur := atoms[i]
		delim := text[prev.Stop:cur.Start]
		lo := i - window
		if lo < 0 {
			lo = 0
		}
		left := atoms[lo:i]
		hi := i + window
		if hi > len(atoms) {
			hi = len(atoms)
		}
		right := atoms[i:hi]
		parts = append(parts, tokenPart{split: &tokenSplit{leftAtoms: left, rightAtoms: right, Delimiter: delim}})
		parts = append(parts, tokenPart{text: atoms[i].Text})
	}
	return parts
}

// TokenTexts returns final token strings (upstream tokenize() without find_substrings offset pass).
func TokenTexts(text string) []string {
	atoms := Atomize(text)
	if len(atoms) == 0 {
		return nil
	}
	parts := tokenSplitterParts(text, atoms)
	return segmentStrings(parts)
}

// TokenSpans returns (start, end) byte offsets for each token, matching upstream find_substrings.
func TokenSpans(text string) [][2]int {
	chunks := TokenTexts(text)
	if len(chunks) == 0 {
		return nil
	}
	out := make([][2]int, 0, len(chunks))
	offset := 0
	for _, chunk := range chunks {
		idx := strings.Index(text[offset:], chunk)
		if idx < 0 {
			panic("tokenize: TokenSpans: chunk not found at expected offset (internal inconsistency)")
		}
		start := offset + idx
		end := start + len(chunk)
		out = append(out, [2]int{start, end})
		offset = end
	}
	return out
}
