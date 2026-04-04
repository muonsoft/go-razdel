package testkit

import "strings"

// ParsePartition splits s into expected segment strings using ASCII '|' as the delimiter.
//
// An empty input yields nil (no segments). A non-empty input uses strings.Split(s, "|"),
// so consecutive pipes yield empty strings (e.g. "a||b" -> "a", "", "b").
func ParsePartition(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "|")
}
