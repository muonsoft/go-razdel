package tokenize_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/muonsoft/go-razdel/internal/testkit"
	"github.com/muonsoft/go-razdel/internal/tokenize"
)

func upstreamTokenTexts(t *testing.T, razdelRoot, input string) []string {
	t.Helper()
	const py = `
import json, sys
from razdel.segmenters.tokenize import tokenize
s = sys.stdin.read()
print(json.dumps([x.text for x in tokenize(s)], ensure_ascii=False))
`
	cmd := exec.Command("python3", "-c", py)
	cmd.Dir = razdelRoot
	cmd.Env = append(os.Environ(), "PYTHONUTF8=1")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Skipf("python3 upstream check skipped: %v", err)
	}
	var texts []string
	if err := json.Unmarshal(out.Bytes(), &texts); err != nil {
		t.Fatalf("decode python output: %v\n%s", err, out.String())
	}
	return texts
}

func TestUpstream_tokenTexts_parity(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	razdelRoot := filepath.Join(root, "third_party", "razdel")
	if _, err := os.Stat(filepath.Join(razdelRoot, "razdel", "segmenters", "tokenize.py")); err != nil {
		t.Skip("third_party/razdel not present:", err)
	}
	cases := []string{
		"mβж",
		"Δσ",
		"1",
		",",
		"...",
		"a b",
		"ab",
		"привет мир",
		"привет, мир",
		"1,5",
		"1/2",
		"что-то",
		"К_тому_же",
		"yahoo!",
		"...?!",
		":-)",
		"***",
	}
	for _, s := range cases {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			want := upstreamTokenTexts(t, razdelRoot, s)
			got := tokenize.TokenTexts(s)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("TokenTexts mismatch\ninput %q\ngo  %#v\npy  %#v", s, got, want)
			}
		})
	}
}

func TestTokenSpans_ascii_matchesPythonCodeUnitIndices(t *testing.T) {
	t.Parallel()
	// Upstream find_substrings uses Python str indices (code units = runes for BMP text).
	// For pure ASCII, those match UTF-8 byte offsets (docs/contracts.md Variant A).
	root := testkit.ModuleRoot(t)
	razdelRoot := filepath.Join(root, "third_party", "razdel")
	if _, err := os.Stat(filepath.Join(razdelRoot, "razdel", "segmenters", "tokenize.py")); err != nil {
		t.Skip("third_party/razdel not present:", err)
	}
	const py = `
import json, sys
from razdel.segmenters.tokenize import tokenize
s = sys.stdin.read()
print(json.dumps([[x.start, x.stop] for x in tokenize(s)], ensure_ascii=False))
`
	input := "a b c"
	cmd := exec.Command("python3", "-c", py)
	cmd.Dir = razdelRoot
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Skipf("python3 upstream check skipped: %v", err)
	}
	var pySpans [][]any
	if err := json.Unmarshal(out.Bytes(), &pySpans); err != nil {
		t.Fatalf("decode: %v", err)
	}
	got := tokenize.TokenSpans(input)
	if !reflect.DeepEqual(got, [][2]int{{0, 1}, {2, 3}, {4, 5}}) {
		t.Fatalf("TokenSpans: got %#v", got)
	}
	for i := range got {
		start := int(pySpans[i][0].(float64))
		end := int(pySpans[i][1].(float64))
		if got[i][0] != start || got[i][1] != end {
			t.Fatalf("span[%d]: py [%d,%d] vs go [%d,%d]", i, start, end, got[i][0], got[i][1])
		}
	}
}
