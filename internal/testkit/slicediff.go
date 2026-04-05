package testkit

import (
	"fmt"
	"strings"
	"testing"
)

// FormatStringSliceMismatch returns a multi-line report when two []string slices differ
// (length, element-wise text, or both). Intended for tokenizer parity failures.
func FormatStringSliceMismatch(want, got []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "length: want %d, got %d\n", len(want), len(got))
	fmt.Fprintf(&b, "want: %#v\n", want)
	fmt.Fprintf(&b, "got:  %#v\n", got)
	n := len(want)
	if len(got) > n {
		n = len(got)
	}
	for i := 0; i < n; i++ {
		switch {
		case i >= len(want):
			fmt.Fprintf(&b, "  [%d]: extra in got: %q\n", i, got[i])
		case i >= len(got):
			fmt.Fprintf(&b, "  [%d]: missing in got (expected %q)\n", i, want[i])
		case want[i] != got[i]:
			fmt.Fprintf(&b, "  [%d]: want %q, got %q\n", i, want[i], got[i])
		}
	}
	return b.String()
}

// AssertStringSliceEqual fails the test with a detailed diff if want and got differ.
func AssertStringSliceEqual(tb testing.TB, want, got []string, context string) {
	tb.Helper()
	if len(want) == len(got) {
		ok := true
		for i := range want {
			if want[i] != got[i] {
				ok = false
				break
			}
		}
		if ok {
			return
		}
	}
	tb.Fatalf("%s\n%s", context, FormatStringSliceMismatch(want, got))
}
