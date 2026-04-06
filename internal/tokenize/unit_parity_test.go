package tokenize_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/fixture"
	"github.com/muonsoft/go-razdel/internal/testkit"
	"github.com/muonsoft/go-razdel/internal/tokenize"
)

// upstreamUnitTokenizeNames align 1:1 with testdata/upstream/unit_tokenize.txt (non-empty,
// non-comment lines in order). See third_party/razdel/razdel/tests/test_tokenize.py UNIT.
var upstreamUnitTokenizeNames = []string{
	"digit_1",
	"dash_chto_to",
	"underscore_K_tomu_zhe",
	"ellipsis",
	"float_1_5",
	"fraction_1_2",
	"guillemet_fill_dot",
	"paren_fill_dot",
	"paren_fill_guillemet_open",
	"smile_colon_paren_paren_paren",
	"smile_colon_paren_fill_comma",
	"mixed_script_m_beta_zhe",
	"delta_sigma",
	"empty_input",
}

func TestUpstream_UNIT_tokenize_partitions(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "testdata", "upstream", "unit_tokenize.txt")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	partitions := fixture.ParsePartitionLines(string(raw))
	if len(partitions) != len(upstreamUnitTokenizeNames) {
		t.Fatalf("unit_tokenize.txt: got %d partitions, want %d (sync upstreamUnitTokenizeNames)",
			len(partitions), len(upstreamUnitTokenizeNames))
	}
	for i, p := range partitions {
		i, p := i, p
		name := upstreamUnitTokenizeNames[i]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			input := p.Text()
			var want []string
			for _, seg := range p.ExpectedSegments() {
				want = append(want, seg.Text)
			}
			got := tokenize.TokenTexts(input)
			ctx := "internal/tokenize.TokenTexts parity\n" +
				"upstream: third_party/razdel/razdel/tests/test_tokenize.py UNIT\n" +
				"fixture: testdata/upstream/unit_tokenize.txt\n" +
				"partition line: " + fixture.FormatPartitionLine(p)
			testkit.AssertStringSliceEqual(t, want, got, ctx)
			if len(input) == 0 {
				if got != nil {
					t.Fatalf("empty input: TokenTexts should be nil slice, got %#v", got)
				}
				if toks := razdel.Tokenize(input); toks != nil {
					t.Fatalf("empty input: Tokenize should be nil, got %#v", toks)
				}
				return
			}
			toks := razdel.Tokenize(input)
			var gotTexts []string
			for _, tok := range toks {
				gotTexts = append(gotTexts, tok.Text)
			}
			testkit.AssertStringSliceEqual(t, want, gotTexts,
				"razdel.Tokenize must match partition expectation\npartition line: "+fixture.FormatPartitionLine(p))
			testkit.AssertTokenOffsetContract(t, input, toks)
		})
	}
}
