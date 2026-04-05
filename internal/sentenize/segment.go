package sentenize

import "strings"

// Segment merges SentSplitterParts using join(s): true means upstream JOIN (merge buffer + delimiter + right).
// parts must follow [string, *SentSplit, string, *SentSplit, ... string]. Buffer is set on each split like upstream.
func Segment(parts []any, join func(*SentSplit) bool) []string {
	if len(parts) == 0 {
		return nil
	}
	buf, ok := parts[0].(string)
	if !ok {
		return nil
	}
	var out []string
	for i := 1; i+1 < len(parts); i += 2 {
		sp, ok := parts[i].(*SentSplit)
		if !ok {
			return nil
		}
		right, ok := parts[i+1].(string)
		if !ok {
			return nil
		}
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
