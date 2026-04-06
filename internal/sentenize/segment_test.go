package sentenize

import (
	"reflect"
	"testing"
)

func TestSegment_trivialParity(t *testing.T) {
	// Expected sentence texts match PYTHONPATH=. python sentenize with trivial rules only
	// (see T012 task / sentenize.py RULES[:4]).
	cases := []struct {
		text string
		want []string
	}{
		{"a.b", []string{"a.b"}},
		{"a. b", []string{"a. b"}},
		{"word.Next", []string{"word.Next"}},
		{"Привет.", []string{"Привет."}},
		{`a."b"`, []string{`a."b"`}},
		{"a.;b", []string{"a.;b"}},
		{"end.!", []string{"end.!"}},
		{"x.:-)y", []string{"x.:-)y"}},
		{"a. B", []string{"a.", "B"}},
		{"Hello. World", []string{"Hello.", "World"}},
		{"x. y", []string{"x. y"}},
		{"no.9", []string{"no.9"}},
		{`a. "Hi"`, []string{"a.", `"Hi"`}},
	}
	for _, tc := range cases {
		parts := SentSplitterParts(tc.text, DefaultWindow)
		raw := Segment(parts, JoinTrivial)
		got := PostStrip(raw)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("text=%q\ngot  %#v\nwant %#v", tc.text, got, tc.want)
		}
	}
}

func TestSegment_singleChunk(t *testing.T) {
	parts := []any{"hello"}
	got := PostStrip(Segment(parts, JoinTrivial))
	if !reflect.DeepEqual(got, []string{"hello"}) {
		t.Fatalf("got %#v", got)
	}
}

func TestByteSpans_roundTrip(t *testing.T) {
	text := "a. B"
	parts := SentSplitterParts(text, DefaultWindow)
	chunks := PostStrip(Segment(parts, JoinTrivial))
	spans := ByteSpans(text, chunks)
	if len(spans) != 2 {
		t.Fatalf("spans=%v", spans)
	}
	if text[spans[0][0]:spans[0][1]] != "a." || text[spans[1][0]:spans[1][1]] != "B" {
		t.Fatalf("spans=%v texts=%q %q", spans, text[spans[0][0]:spans[0][1]], text[spans[1][0]:spans[1][1]])
	}
}

func TestByteSpans_panicsWhenChunkNotFound(t *testing.T) {
	text := "ab"
	chunks := []string{"a", "missing"}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when chunk is not found")
		}
	}()
	_ = ByteSpans(text, chunks)
}
