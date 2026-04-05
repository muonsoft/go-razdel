package sentenize

import "strings"

// ByteSpans maps stripped chunks to byte offsets in text (upstream find_substrings, UTF-8 bytes).
func ByteSpans(text string, chunks []string) [][2]int {
	offset := 0
	out := make([][2]int, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk == "" {
			out = append(out, [2]int{offset, offset})
			continue
		}
		idx := strings.Index(text[offset:], chunk)
		if idx < 0 {
			idx = strings.Index(text, chunk)
			if idx < 0 {
				continue
			}
		}
		start := offset + idx
		end := start + len(chunk)
		out = append(out, [2]int{start, end})
		offset = end
	}
	return out
}
