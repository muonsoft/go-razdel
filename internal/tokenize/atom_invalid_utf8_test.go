package tokenize_test

import (
	"testing"

	"github.com/muonsoft/go-razdel/internal/tokenize"
)

func TestAtomize_invalidUTF8_eachByteOther(t *testing.T) {
	t.Parallel()
	src := string([]byte{'a', 0xff, 'b'})
	got := tokenize.Atomize(src)
	if len(got) != 3 {
		t.Fatalf("len=%d want 3: %#v", len(got), got)
	}
	want := []tokenize.Atom{
		{Start: 0, Stop: 1, Type: tokenize.LAT, Text: "a"},
		{Start: 1, Stop: 2, Type: tokenize.OTHER, Text: "\xff"},
		{Start: 2, Stop: 3, Type: tokenize.LAT, Text: "b"},
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("atom[%d]: got %#v want %#v", i, got[i], want[i])
		}
	}
}

func TestAtomize_invalidUTF8_eachInvalidByteSeparateAtom(t *testing.T) {
	t.Parallel()
	src := string([]byte{0xff, 0xfe, 0x80})
	got := tokenize.Atomize(src)
	if len(got) != 3 {
		t.Fatalf("len=%d want 3: %#v", len(got), got)
	}
	for i, a := range got {
		if a.Type != tokenize.OTHER || a.Start != i || a.Stop != i+1 || a.Text != src[i:i+1] {
			t.Fatalf("atom[%d]: got %#v", i, a)
		}
	}
}
