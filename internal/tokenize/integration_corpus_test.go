package tokenize_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/muonsoft/go-razdel/internal/fixture"
	"github.com/muonsoft/go-razdel/internal/testkit"
	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// envTokenizeIntegrationFull enables the full upstream tokens.txt integration run
// (third_party/razdel/razdel/tests/data/tokens.txt). Default CI and `go test ./...`
// use only the deterministic quick sample (testdata/upstream/quick_tokens_sample.txt).
const envTokenizeIntegrationFull = "RAZDEL_TOKENIZE_INTEGRATION_FULL"

func TestIntegration_tokenize_quick_corpus(t *testing.T) {
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "testdata", "upstream", "quick_tokens_sample.txt")
	runTokenizeCorpus(t, path, "quick", true)
}

func TestIntegration_tokenize_full_corpus(t *testing.T) {
	if os.Getenv(envTokenizeIntegrationFull) != "1" {
		t.Skipf("full corpus: set %s=1 (see testdata/upstream/README.md)", envTokenizeIntegrationFull)
	}
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "third_party", "razdel", "razdel", "tests", "data", "tokens.txt")
	if _, err := os.Stat(path); err != nil {
		t.Skip("tokens.txt not available:", err)
	}
	runTokenizeCorpus(t, path, "full", false)
}

func runTokenizeCorpus(t *testing.T, path, mode string, subtests bool) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parts := fixture.ParsePartitionLines(string(raw))
	n := len(parts)
	if n == 0 {
		t.Fatal("corpus: no partition lines parsed from", path)
	}
	t.Logf("tokenizer integration (%s): %d cases from %s", mode, n, path)

	start := time.Now()
	if subtests {
		for i, p := range parts {
			i, p := i, p
			t.Run(strconv.Itoa(i+1), func(t *testing.T) {
				t.Parallel()
				assertPartitionTokenize(t, mode, path, i+1, n, p)
			})
		}
		t.Logf("tokenizer integration (%s): finished in %s", mode, time.Since(start).Round(time.Millisecond))
		return
	}

	var mismatches []corpusMismatch
	var skippedDrift int
	for i, p := range parts {
		lineKey := fixture.PartitionLineKey(p)
		if _, drift := upstreamTokensTxtPartitionDrift[lineKey]; drift {
			skippedDrift++
			continue
		}
		input := p.Text()
		var want []string
		for _, seg := range p.ExpectedSegments() {
			want = append(want, seg.Text)
		}
		got := tokenize.TokenTexts(input)
		if slices.Equal(want, got) {
			continue
		}
		mismatches = append(mismatches, corpusMismatch{
			index: i + 1,
			total: n,
			part:  p,
			want:  want,
			got:   got,
		})
	}
	elapsed := time.Since(start).Round(time.Millisecond)
	if skippedDrift > 0 {
		t.Logf("tokenizer integration (%s): skipped %d lines listed in upstream_tokens_txt_drift_test.go (tokens.txt etalon ≠ tokenize.py on same text in pinned submodule)",
			mode, skippedDrift)
	}
	if len(mismatches) == 0 {
		t.Logf("tokenizer integration (%s): all %d checked partition cases matched in %s", mode, n-skippedDrift, elapsed)
		if mode == "full" {
			if skippedDrift != len(upstreamTokensTxtPartitionDrift) {
				t.Fatalf("upstream drift skip list has %d entries but corpus contained only %d of them (update upstream_tokens_txt_drift_test.go or submodule)",
					len(upstreamTokensTxtPartitionDrift), skippedDrift)
			}
			documentUpstreamDriftSkips(t)
		}
		return
	}
	t.Logf("tokenizer integration (%s): %d mismatches out of %d partition cases (elapsed %s)",
		mode, len(mismatches), n, elapsed)
	const maxLogged = 20
	for j, m := range mismatches {
		if j >= maxLogged {
			t.Logf("... omitting %d more mismatches (see Fatalf summary)", len(mismatches)-maxLogged)
			break
		}
		ctx := formatCorpusMismatchContext(mode, path, m.index, m.total, m.part)
		t.Logf("%s\n%s", ctx, testkit.FormatStringSliceMismatch(m.want, m.got))
	}
	t.Fatalf("corpus: %d partition lines where tokens.txt etalon != tokenize.TokenTexts (pinned submodule: razdel.segmenters.tokenize disagrees on the same lines; fixture drift, see testdata/upstream/README.md)",
		len(mismatches))
}

type corpusMismatch struct {
	index, total int
	part         fixture.Partition
	want, got    []string
}

func assertPartitionTokenize(t *testing.T, mode, path string, index, total int, p fixture.Partition) {
	t.Helper()
	input := p.Text()
	var want []string
	for _, seg := range p.ExpectedSegments() {
		want = append(want, seg.Text)
	}
	got := tokenize.TokenTexts(input)
	ctx := formatCorpusMismatchContext(mode, path, index, total, p)
	testkit.AssertStringSliceEqual(t, want, got, ctx)
}

func formatCorpusMismatchContext(mode, path string, index, total int, p fixture.Partition) string {
	return fmt.Sprintf(
		"tokenize corpus parity (%s)\n"+
			"file: %s\n"+
			"case: %d/%d\n"+
			"upstream: third_party/razdel/razdel/tests/test_tokenize.py test_int / int_tests\n"+
			"partition line: %s",
		mode, path, index, total, fixture.FormatPartitionLine(p),
	)
}

// documentUpstreamDriftSkips emits one skipped subtest per known drift line so -v output and
// reports show SKIP with an explicit reason (tokens.txt row vs tokenize.py in pinned upstream).
func documentUpstreamDriftSkips(t *testing.T) {
	t.Helper()
	keys := make([]string, 0, len(upstreamTokensTxtPartitionDrift))
	for k := range upstreamTokensTxtPartitionDrift {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for i, line := range keys {
		i, line := i, line
		t.Run(fmt.Sprintf("upstream_tokens_txt_drift_%02d", i+1), func(t *testing.T) {
			t.Skipf("tokens.txt partition %q: file etalon != tokenize.py on same text in pinned submodule; Go matches Python (see upstream_tokens_txt_drift_test.go)", line)
		})
	}
}
