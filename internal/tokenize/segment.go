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
	atoms     []Atom
	index     int // index of right atom in atoms slice (left atom is index-1)
	Delimiter string
	buffer    string // current buffered token text in segment pass
}

func (s *tokenSplit) left1() *Atom  { return &s.atoms[s.index-1] }
func (s *tokenSplit) left2() *Atom  { return &s.atoms[s.index-2] }
func (s *tokenSplit) right1() *Atom { return &s.atoms[s.index] }
func (s *tokenSplit) right2() *Atom { return &s.atoms[s.index+1] }

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
		if s.index < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case dashDelimiterPiece(s.Right()):
		if s.index+1 >= len(s.atoms) {
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
		if s.index < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case underscoreDelimiterPiece(s.Right()):
		if s.index+1 >= len(s.atoms) {
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
		if s.index < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case floatDelimiterPiece(s.Right()):
		if s.index+1 >= len(s.atoms) {
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
		if s.index < 2 {
			return 0
		}
		left, right = s.left2(), s.right1()
	case fractionDelimiterPiece(s.Right()):
		if s.index+1 >= len(s.atoms) {
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
	if s.index+1 < len(s.atoms) && smileRE.MatchString(combined+s.right2().Text) {
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

var tokenJoinRules = [...]func(*tokenSplit) int{
	dashRule2112,
	underscoreRule2112,
	floatRule2112,
	fractionRule2112,
	punctRule,
	otherRule,
	yahooRule,
}

func joinSplit(s *tokenSplit) bool {
	for _, rule := range tokenJoinRules {
		if a := rule(s); a != 0 {
			return a == join
		}
	}
	return false
}

func tokenSpansFromAtoms(text string, atoms []Atom) [][2]int {
	if len(atoms) == 0 {
		return nil
	}
	start := atoms[0].Start
	end := atoms[0].Stop
	out := make([][2]int, 0, len(atoms))
	for i := 1; i < len(atoms); i++ {
		prev := atoms[i-1]
		cur := atoms[i]
		split := tokenSplit{
			atoms:     atoms,
			index:     i,
			Delimiter: text[prev.Stop:cur.Start],
			buffer:    text[start:end],
		}
		if shouldJoin(&split) {
			end = cur.Stop
			continue
		}
		out = append(out, [2]int{start, end})
		start = cur.Start
		end = cur.Stop
	}
	out = append(out, [2]int{start, end})
	return out
}

// TokenTexts returns final token strings (upstream tokenize() without find_substrings offset pass).
func TokenTexts(text string) []string {
	atoms := Atomize(text)
	if len(atoms) == 0 {
		return nil
	}
	spans := tokenSpansFromAtoms(text, atoms)
	out := make([]string, len(spans))
	for i, span := range spans {
		out[i] = text[span[0]:span[1]]
	}
	return out
}

// TokenSpans returns (start, end) byte offsets for each token, matching upstream find_substrings.
func TokenSpans(text string) [][2]int {
	atoms := Atomize(text)
	if len(atoms) == 0 {
		return nil
	}
	return tokenSpansFromAtoms(text, atoms)
}
