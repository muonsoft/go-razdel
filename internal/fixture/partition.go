package fixture

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// FillPattern matches upstream partition.py: chunks that are only whitespace are "fill" gaps.
var fillPattern = regexp.MustCompile(`^\s*$`)

// IsFill reports whether chunk is treated as a gap between segments (upstream Partition.is_fill).
func IsFill(chunk string) bool {
	return fillPattern.MatchString(chunk)
}

// Partition is one upstream-style partition line: chunks split by ASCII '|' (see razdel/tests/partition.py).
type Partition struct {
	Chunks []string
}

// Text is the reconstructed input string (concatenation of all chunks).
func (p Partition) Text() string {
	return strings.Join(p.Chunks, "")
}

// ExpectedSegment is a non-fill substring with rune offsets into Text(), matching upstream Substring
// start/stop semantics (Python len(str) counts Unicode code points).
type ExpectedSegment struct {
	StartRune int
	EndRune   int
	Text      string
}

// ExpectedSegments returns non-fill chunks with half-open [StartRune, EndRune) rune indices.
func (p Partition) ExpectedSegments() []ExpectedSegment {
	var out []ExpectedSegment
	runePos := 0
	for _, chunk := range p.Chunks {
		nRunes := utf8.RuneCountInString(chunk)
		if !IsFill(chunk) {
			out = append(out, ExpectedSegment{
				StartRune: runePos,
				EndRune:   runePos + nRunes,
				Text:      chunk,
			})
		}
		runePos += nRunes
	}
	return out
}

// EmptyPartitionMarker is a single-line sentinel in fixture files for the upstream UNIT case of
// an empty partition line (text "", no segments). A raw empty line cannot be represented in
// newline-delimited files; see testdata/upstream/README.md.
const EmptyPartitionMarker = "#empty"

// ParsePartitionLine parses one line like upstream parse_partition: split on '|' with no escape processing.
// For line "" (empty string), Python ''.split('|') yields [''].
func ParsePartitionLine(line string) Partition {
	return Partition{Chunks: strings.Split(line, "|")}
}

// ParsePartitionLines parses newline-separated partition records. Blank lines are skipped.
// A line exactly equal to EmptyPartitionMarker denotes an empty partition (one fill chunk).
func ParsePartitionLines(content string) []Partition {
	content = strings.TrimSuffix(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	lines := strings.Split(content, "\n")
	var out []Partition
	for _, line := range lines {
		if line == "" {
			continue
		}
		if line == EmptyPartitionMarker {
			out = append(out, Partition{Chunks: []string{""}})
			continue
		}
		out = append(out, ParsePartitionLine(line))
	}
	return out
}
