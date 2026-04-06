package razdel_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muonsoft/go-razdel/internal/fixture"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

func TestUpstreamQuickSamples_loadAndParse(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	for _, name := range []string{
		"testdata/upstream/quick_tokens_sample.txt",
		"testdata/upstream/quick_sents_sample.txt",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(root, filepath.FromSlash(name))
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			parts := fixture.ParsePartitionLines(string(raw))
			if len(parts) == 0 {
				t.Fatal("expected at least one partition")
			}
			for i, p := range parts {
				text := p.Text()
				for _, seg := range p.ExpectedSegments() {
					runes := []rune(text)
					if seg.StartRune < 0 || seg.EndRune > len(runes) || seg.StartRune > seg.EndRune {
						t.Fatalf("partition %d: bad rune span [%d,%d)", i, seg.StartRune, seg.EndRune)
					}
					got := string(runes[seg.StartRune:seg.EndRune])
					if got != seg.Text {
						t.Fatalf("partition %d: rune slice %q != segment %q", i, got, seg.Text)
					}
				}
			}
		})
	}
}

func TestUpstreamMETA_documentsSubmodule(t *testing.T) {
	t.Parallel()
	path := filepath.Join(testkit.ModuleRoot(t), "testdata", "upstream", "META.txt")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(raw)
	if !strings.Contains(content, "source_submodule_commit=") {
		t.Fatal("META.txt should list source_submodule_commit")
	}
	if !strings.Contains(content, "sample_algorithm=") {
		t.Fatal("META.txt should list sample_algorithm")
	}
}
