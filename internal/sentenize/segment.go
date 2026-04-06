package sentenize

import "strings"

// Segment merges SentSplitterParts using join(s): true means upstream JOIN (merge buffer + delimiter + right).
// parts must follow [text, split, text, split, ... text]. Buffer is set on each split like upstream.
func Segment(parts []SentPart, join func(*SentSplit) bool) []string {
	if len(parts) == 0 {
		return nil
	}
	buf := parts[0].Text
	var out []string
	for i := 1; i+1 < len(parts); i += 2 {
		sp := parts[i].Split
		right := parts[i+1].Text
		bufCopy := buf
		sp.Buffer = &bufCopy
		if join(sp) {
			buf = buf + sp.Delimiter + right
		} else {
			out = append(out, buf+sp.Delimiter)
			buf = right
		}
	}
	out = append(out, buf)
	return out
}

// PostStrip applies upstream SentSegmenter.post: strings.TrimSpace on each chunk.
func PostStrip(chunks []string) []string {
	out := make([]string, len(chunks))
	for i, c := range chunks {
		out[i] = strings.TrimSpace(c)
	}
	return out
}
