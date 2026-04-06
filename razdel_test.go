package razdel_test

import (
	"strings"
	"testing"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

func TestTokenize_emptyAndCyrillic(t *testing.T) {
	t.Parallel()
	empty := razdel.Tokenize("")
	if len(empty) != 0 {
		t.Fatalf("empty input: want len 0, got %d", len(empty))
	}
	testkit.AssertTokenOffsetContract(t, "", empty)

	src := "привет"
	got := razdel.Tokenize(src)
	testkit.AssertTokenTextsEqual(t, src, got, []string{"привет"})

	testkit.AssertTokenTextsEqual(t, "a b", razdel.Tokenize("a b"), []string{"a", "b"})
	testkit.AssertTokenTextsEqual(t, "привет, мир", razdel.Tokenize("привет, мир"), []string{"привет", ",", "мир"})
}

func TestSentenize_emptyWhitespaceAndTrivial(t *testing.T) {
	t.Parallel()
	empty := razdel.Sentenize("")
	if len(empty) != 0 {
		t.Fatalf("empty input: want len 0, got %d", len(empty))
	}
	testkit.AssertSentenceOffsetContract(t, "", empty)

	ws := razdel.Sentenize("  \n\t  ")
	if len(ws) != 0 {
		t.Fatalf("whitespace-only: want len 0, got %d", len(ws))
	}
	testkit.AssertSentenceOffsetContract(t, "  \n\t  ", ws)

	// strings.TrimSpace strips the same class of Unicode spaces as Python str.strip() for common cases.
	for _, s := range []string{"\u00a0", "\u2003\u2002\u2009", "\u1680\u2000"} {
		if strings.TrimSpace(s) != "" {
			t.Fatalf("precondition: TrimSpace(%q) should be empty", s)
		}
		got := razdel.Sentenize(s)
		if len(got) != 0 {
			t.Fatalf("unicode whitespace-only %q: want len 0, got %d (%#v)", s, len(got), got)
		}
		testkit.AssertSentenceOffsetContract(t, s, got)
	}

	src := "Привет."
	got := razdel.Sentenize(src)
	if len(got) != 1 || got[0].Text != src {
		t.Fatalf("got %#v want single sentence %q", got, src)
	}
	testkit.AssertSentenceOffsetContract(t, src, got)

	// Upstream list_item: single-letter "a" is in BULLET_CHARS, so "a. B" stays one sentence.
	joined := razdel.Sentenize("a. B")
	if len(joined) != 1 || joined[0].Text != "a. B" {
		t.Fatalf("got %#v want single sentence %q (upstream sentenize)", joined, "a. B")
	}
	testkit.AssertSentenceOffsetContract(t, "a. B", joined)

	// Post-step strip: offsets must point at trimmed chunks inside the original string.
	padded := "  Одно. Два.  "
	padGot := razdel.Sentenize(padded)
	testkit.AssertSentenceTextsEqual(t, padded, padGot, []string{"Одно.", "Два."})
}
