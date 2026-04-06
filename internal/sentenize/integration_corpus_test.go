package sentenize_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/fixture"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

// envSentenizeIntegrationFull enables the full upstream sents.txt integration run
// (third_party/razdel/razdel/tests/data/sents.txt). Default CI and `go test ./...`
// use only the deterministic quick sample (testdata/upstream/quick_sents_sample.txt).
const envSentenizeIntegrationFull = "RAZDEL_SENTENIZE_INTEGRATION_FULL"

func TestIntegration_sentenize_quick_corpus(t *testing.T) {
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "testdata", "upstream", "quick_sents_sample.txt")
	runSentenizeCorpus(t, path, "quick", true)
}

func TestIntegration_sentenize_full_corpus(t *testing.T) {
	if os.Getenv(envSentenizeIntegrationFull) != "1" {
		t.Skipf("full corpus: set %s=1 (see testdata/upstream/README.md)", envSentenizeIntegrationFull)
	}
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "third_party", "razdel", "razdel", "tests", "data", "sents.txt")
	if _, err := os.Stat(path); err != nil {
		t.Skip("sents.txt not available:", err)
	}
	runSentenizeCorpus(t, path, "full", false)
}

func runSentenizeCorpus(t *testing.T, path, mode string, subtests bool) {
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
	t.Logf("sentenizer integration (%s): %d cases from %s", mode, n, path)

	start := time.Now()
	if subtests {
		for i, p := range parts {
			i, p := i, p
			t.Run(strconv.Itoa(i+1), func(t *testing.T) {
				t.Parallel()
				assertPartitionSentenize(t, mode, path, i+1, n, p)
			})
		}
		t.Logf("sentenizer integration (%s): finished in %s", mode, time.Since(start).Round(time.Millisecond))
		return
	}

	var mismatches []corpusMismatch
	var skippedDrift int
	for i, p := range parts {
		lineKey := fixture.PartitionLineKey(p)
		if _, drift := upstreamSentsTxtPartitionDrift[lineKey]; drift {
			skippedDrift++
			continue
		}
		input := p.Text()
		var want []string
		for _, seg := range p.ExpectedSegments() {
			want = append(want, seg.Text)
		}
		got := sentenceTexts(t, input)
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
		t.Logf("sentenizer integration (%s): skipped %d lines listed in upstream_sents_txt_drift_test.go (sents.txt etalon ≠ sentenize.py on same text in pinned submodule)",
			mode, skippedDrift)
	}
	if len(mismatches) == 0 {
		t.Logf("sentenizer integration (%s): all %d checked partition cases matched in %s", mode, n-skippedDrift, elapsed)
		if mode == "full" {
			if skippedDrift != len(upstreamSentsTxtPartitionDrift) {
				t.Fatalf("upstream drift skip list has %d entries but corpus contained only %d of them (update upstream_sents_txt_drift_test.go or submodule)",
					len(upstreamSentsTxtPartitionDrift), skippedDrift)
			}
			documentUpstreamDriftSkips(t)
		}
		return
	}
	t.Logf("sentenizer integration (%s): %d mismatches out of %d partition cases (elapsed %s)",
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
	t.Fatalf("corpus: %d partition lines where sents.txt etalon != razdel.Sentenize (pinned submodule: razdel.segmenters.sentenize disagrees on the same lines; fixture drift, see testdata/upstream/README.md)",
		len(mismatches))
}

type corpusMismatch struct {
	index, total int
	part         fixture.Partition
	want, got    []string
}

func sentenceTexts(t *testing.T, text string) []string {
	t.Helper()
	sents := razdel.Sentenize(text)
	testkit.AssertSentenceOffsetContract(t, text, sents)
	out := make([]string, len(sents))
	for i, s := range sents {
		out[i] = s.Text
	}
	return out
}

func assertPartitionSentenize(t *testing.T, mode, path string, index, total int, p fixture.Partition) {
	t.Helper()
	input := p.Text()
	var want []string
	for _, seg := range p.ExpectedSegments() {
		want = append(want, seg.Text)
	}
	got := sentenceTexts(t, input)
	ctx := formatCorpusMismatchContext(mode, path, index, total, p)
	testkit.AssertStringSliceEqual(t, want, got, ctx)
}

func formatCorpusMismatchContext(mode, path string, index, total int, p fixture.Partition) string {
	return fmt.Sprintf(
		"sentenize corpus parity (%s)\n"+
			"file: %s\n"+
			"case: %d/%d\n"+
			"upstream: third_party/razdel/razdel/tests/test_sentenize.py test_int / int_tests\n"+
			"partition line: %s",
		mode, path, index, total, fixture.FormatPartitionLine(p),
	)
}

// documentUpstreamDriftSkips emits one skipped subtest per known drift line so -v output and
// reports show SKIP with an explicit reason (sents.txt row vs sentenize.py in pinned upstream).
func documentUpstreamDriftSkips(t *testing.T) {
	t.Helper()
	keys := make([]string, 0, len(upstreamSentsTxtPartitionDrift))
	for k := range upstreamSentsTxtPartitionDrift {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for i, line := range keys {
		i, line := i, line
		t.Run(fmt.Sprintf("upstream_sents_txt_drift_%02d", i+1), func(t *testing.T) {
			t.Skipf("sents.txt partition %q: file etalon != sentenize.py on same text in pinned submodule; Go matches Python (see upstream_sents_txt_drift_test.go)", line)
		})
	}
}
