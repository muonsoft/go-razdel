package razdel

import "testing"

func TestTokenize_emptyAndStub(t *testing.T) {
	t.Parallel()
	if got := Tokenize(""); len(got) != 0 {
		t.Fatalf("empty input: want len 0, got %d", len(got))
	}
	if got := Tokenize("привет"); len(got) != 0 {
		t.Fatalf("stub: non-empty input must still return empty, got len %d", len(got))
	}
}

func TestSentenize_emptyAndStub(t *testing.T) {
	t.Parallel()
	if got := Sentenize(""); len(got) != 0 {
		t.Fatalf("empty input: want len 0, got %d", len(got))
	}
	if got := Sentenize("Привет."); len(got) != 0 {
		t.Fatalf("stub: non-empty input must still return empty, got len %d", len(got))
	}
}
