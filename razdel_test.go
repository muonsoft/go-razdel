package razdel_test

import (
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

	src := "Привет."
	got := razdel.Sentenize(src)
	if len(got) != 1 || got[0].Text != src {
		t.Fatalf("got %#v want single sentence %q", got, src)
	}
	testkit.AssertSentenceOffsetContract(t, src, got)

	split := razdel.Sentenize("a. B")
	if len(split) != 2 || split[0].Text != "a." || split[1].Text != "B" {
		t.Fatalf("got %#v want [a. B]", split)
	}
	testkit.AssertSentenceOffsetContract(t, "a. B", split)
}
