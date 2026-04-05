package sentenize

import (
	"reflect"
	"testing"
)

func TestSentSplitterParts_whitespaceOnly(t *testing.T) {
	for _, s := range []string{"", "   ", "\n\t  "} {
		if p := SentSplitterParts(s, DefaultWindow); p != nil {
			t.Fatalf("whitespace %q: got %#v, want nil", s, p)
		}
	}
}

func TestSentSplitterParts_upstreamStream(t *testing.T) {
	// Expected slices mirror third_party/razdel/razdel/segmenters/sentenize.py SentSplitter().__call__.
	cases := []struct {
		text string
		want []any
	}{
		{
			"a.b",
			[]any{"a", &SentSplit{Left: "a", Delimiter: ".", Right: "b"}, "b"},
		},
		{
			":-)x",
			[]any{"", &SentSplit{Left: "", Delimiter: ":-)", Right: "x"}, "x"},
		},
		{
			"  a.b  ",
			[]any{"  a", &SentSplit{Left: "  a", Delimiter: ".", Right: "b  "}, "b  "},
		},
		{
			"привет.",
			[]any{"привет", &SentSplit{Left: "привет", Delimiter: ".", Right: ""}, ""},
		},
		{
			"a„b",
			[]any{"a", &SentSplit{Left: "a", Delimiter: "„", Right: "b"}, "b"},
		},
		{
			"word;x",
			[]any{"word", &SentSplit{Left: "word", Delimiter: ";", Right: "x"}, "x"},
		},
		{
			"(x)",
			[]any{"(x", &SentSplit{Left: "(x", Delimiter: ")", Right: ""}, ""},
		},
	}
	for _, tc := range cases {
		got := SentSplitterParts(tc.text, DefaultWindow)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("text=%q\n got %#v\nwant %#v", tc.text, got, tc.want)
		}
	}
}

func TestSentSplit_tokenAccessors(t *testing.T) {
	cases := []struct {
		sp   SentSplit
		ltok string
		rtok string
		rsp  bool
		lss  bool
		rw   string
	}{
		{SentSplit{Left: "left ", Delimiter: ".", Right: " right"}, "left", "right", true, true, "right"},
		{SentSplit{Left: "word", Delimiter: "!", Right: "Next"}, "word", "Next", false, false, "Next"},
		{SentSplit{Left: "", Delimiter: `"`, Right: "quote"}, "", "quote", false, false, "quote"},
		{SentSplit{Left: "x (", Delimiter: ")", Right: " y"}, "(", "y", true, false, "y"},
		{SentSplit{Left: "see ", Delimiter: ":", Right: "-)"}, "see", "-", false, true, ""},
	}
	for _, tc := range cases {
		sp := tc.sp
		if g := sp.LeftToken(); g != tc.ltok {
			t.Errorf("%+v LeftToken()=%q want %q", sp, g, tc.ltok)
		}
		if g := sp.RightToken(); g != tc.rtok {
			t.Errorf("%+v RightToken()=%q want %q", sp, g, tc.rtok)
		}
		if g := sp.RightSpacePrefix(); g != tc.rsp {
			t.Errorf("%+v RightSpacePrefix()=%v want %v", sp, g, tc.rsp)
		}
		if g := sp.LeftSpaceSuffix(); g != tc.lss {
			t.Errorf("%+v LeftSpaceSuffix()=%v want %v", sp, g, tc.lss)
		}
		if g := sp.RightWord(); g != tc.rw {
			t.Errorf("%+v RightWord()=%q want %q", sp, g, tc.rw)
		}
	}
}

func TestSentSplit_sokrAccessors(t *testing.T) {
	sp := SentSplit{Left: "a. b", Delimiter: ".", Right: "c"}
	a, b, ok := sp.LeftPairSokr()
	if !ok || a != "a" || b != "b" {
		t.Fatalf("LeftPairSokr: (%q,%q) ok=%v", a, b, ok)
	}
	sp = SentSplit{Left: "num 12 - word", Delimiter: ".", Right: "x"}
	w, ok := sp.LeftIntSokr()
	if !ok || w != "word" {
		t.Fatalf("LeftIntSokr: %q ok=%v", w, ok)
	}
	sp = SentSplit{Left: "x ", Delimiter: ".", Right: "y"}
	if _, _, ok := sp.LeftPairSokr(); ok {
		t.Fatal("LeftPairSokr: expected no match")
	}
	if _, ok := sp.LeftIntSokr(); ok {
		t.Fatal("LeftIntSokr: expected no match")
	}
}

func TestSentSplit_bufferTokens(t *testing.T) {
	sp := SentSplit{Left: "a", Delimiter: ".", Right: "b"}
	if sp.BufferTokens() != nil || sp.BufferFirstToken() != "" {
		t.Fatal("unset buffer should not expose tokens")
	}
	buf := "a. b, c"
	sp.Buffer = &buf
	want := []string{"a", ".", "b", ",", "c"}
	if g := sp.BufferTokens(); !reflect.DeepEqual(g, want) {
		t.Fatalf("BufferTokens: %#v want %#v", g, want)
	}
	if g := sp.BufferFirstToken(); g != "a" {
		t.Fatalf("BufferFirstToken: %q", g)
	}
}

func TestSentSplitterParts_noDelimiter(t *testing.T) {
	got := SentSplitterParts("hello", DefaultWindow)
	want := []any{"hello"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
