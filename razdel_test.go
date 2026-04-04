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
}

func TestSentenize_emptyAndStub(t *testing.T) {
	t.Parallel()
	empty := razdel.Sentenize("")
	if len(empty) != 0 {
		t.Fatalf("empty input: want len 0, got %d", len(empty))
	}
	testkit.AssertSentenceOffsetContract(t, "", empty)

	got := razdel.Sentenize("Привет.")
	if len(got) != 0 {
		t.Fatalf("stub: non-empty input must still return empty, got len %d", len(got))
	}
	testkit.AssertSentenceOffsetContract(t, "Привет.", got)
}
