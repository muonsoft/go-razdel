package tokenize

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Atom kinds mirror upstream razdel.segmenters.tokenize (RU, LAT, INT, PUNCT, OTHER).
const (
	RU    = "RU"
	LAT   = "LAT"
	INT   = "INT"
	PUNCT = "PUNCT"
	OTHER = "OTHER"
)

// Atom is one lexical atom from upstream Atom (start/stop as byte offsets, half-open).
type Atom struct {
	Start int
	Stop  int
	Type  string
	Text  string
}

// Normal is lowercase text (upstream atom.normal).
func (a Atom) Normal() string {
	return strings.ToLower(a.Text)
}

// punctLookup is pre-computed from PunctSet for O(1) rune membership in Atomize.
var punctLookup = makePunctLookup()

// Atomize scans text into atoms using the same classification order as upstream ATOM regex.
func Atomize(text string) []Atom {
	out := make([]Atom, 0, len(text)/4+1)
	for i := 0; i < len(text); {
		r, w := utf8.DecodeRuneInString(text[i:])
		// Upstream ATOM.finditer never yields atoms for whitespace gaps (only \S via OTHER).
		if unicode.IsSpace(r) {
			i += w
			continue
		}
		if r == utf8.RuneError && w == 1 {
			// Invalid UTF-8: treat as single-byte OTHER (upstream passes through \S).
			out = append(out, Atom{Start: i, Stop: i + 1, Type: OTHER, Text: text[i : i+1]})
			i++
			continue
		}
		start := i
		// RU: [а-яё]+ case-insensitive
		if isCyrillicLetter(r) {
			j := i + w
			for j < len(text) {
				r2, w2 := utf8.DecodeRuneInString(text[j:])
				if r2 == utf8.RuneError && w2 == 1 {
					break
				}
				if !isCyrillicLetter(r2) {
					break
				}
				j += w2
			}
			out = append(out, Atom{Start: start, Stop: j, Type: RU, Text: text[start:j]})
			i = j
			continue
		}
		// LAT: [a-z]+ case-insensitive (ASCII only, matches Python with re.I on [a-z]).
		if isASCIILetter(r) {
			j := i + w
			for j < len(text) {
				r2, w2 := utf8.DecodeRuneInString(text[j:])
				if !isASCIILetter(r2) {
					break
				}
				j += w2
			}
			out = append(out, Atom{Start: start, Stop: j, Type: LAT, Text: text[start:j]})
			i = j
			continue
		}
		// INT: \d+ with Unicode digits (Python re.U).
		if unicode.IsDigit(r) {
			j := i + w
			for j < len(text) {
				r2, w2 := utf8.DecodeRuneInString(text[j:])
				if !unicode.IsDigit(r2) {
					break
				}
				j += w2
			}
			out = append(out, Atom{Start: start, Stop: j, Type: INT, Text: text[start:j]})
			i = j
			continue
		}
		if punctLookup[r] {
			out = append(out, Atom{Start: start, Stop: i + w, Type: PUNCT, Text: text[start : i+w]})
			i += w
			continue
		}
		// OTHER: \S — single non-space code point (upstream); space is skipped, never an atom.
		out = append(out, Atom{Start: start, Stop: i + w, Type: OTHER, Text: text[start : i+w]})
		i += w
	}
	return out
}

func makePunctLookup() map[rune]bool {
	m := make(map[rune]bool, len(PunctSet))
	for _, r := range PunctSet {
		m[r] = true
	}
	return m
}

func isCyrillicLetter(r rune) bool {
	if r == utf8.RuneError {
		return false
	}
	lr := unicode.ToLower(r)
	return (lr >= 'а' && lr <= 'я') || lr == 'ё'
}

func isASCIILetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
