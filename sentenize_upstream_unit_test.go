package razdel_test

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/fixture"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

// Upstream third_party/razdel/razdel/tests/test_sentenize.py UNIT (full list, 34 cases).
//
//go:embed testdata/upstream/unit_sentenize.txt
var unitSentenizeUpstream string

func TestSentenize_upstreamUnitParity(t *testing.T) {
	t.Parallel()
	parts := fixture.ParsePartitionLines(unitSentenizeUpstream)
	if len(parts) != 34 {
		t.Fatalf("fixture: want 34 UNIT partitions (test_sentenize.py), got %d", len(parts))
	}
	for i, p := range parts {
		src := p.Text()
		segs := p.ExpectedSegments()
		want := make([]string, len(segs))
		for j := range segs {
			want[j] = segs[j].Text
		}
		t.Run(fmt.Sprintf("L%02d", i+1), func(t *testing.T) {
			t.Parallel()
			got := razdel.Sentenize(src)
			testkit.AssertSentenceTextsEqual(t, src, got, want)
		})
	}
}
