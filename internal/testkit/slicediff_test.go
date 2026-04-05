package testkit

import (
	"strings"
	"testing"
)

func TestFormatStringSliceMismatch(t *testing.T) {
	t.Parallel()
	s := FormatStringSliceMismatch([]string{"a", "b"}, []string{"a", "c"})
	if !strings.Contains(s, `[1]: want "b", got "c"`) {
		t.Fatalf("expected index 1 mismatch, got:\n%s", s)
	}
	s2 := FormatStringSliceMismatch([]string{"x"}, []string{"x", "y"})
	if !strings.Contains(s2, "extra in got") {
		t.Fatalf("expected extra element note, got:\n%s", s2)
	}
	s3 := FormatStringSliceMismatch([]string{"a", "b"}, []string{"a"})
	if !strings.Contains(s3, "missing in got") {
		t.Fatalf("expected missing element note, got:\n%s", s3)
	}
}
