// Command genupstreamfixtures regenerates sampled corpus fixtures under testdata/upstream from
// third_party/razdel. Run from the repository root: go run ./tools/genupstreamfixtures
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/muonsoft/go-razdel/internal/fixture"
)

const (
	upstreamSubmoduleSHA = "668dbe191a5cfd94bebf9155e2ffa5f94ff3fe33"
	sampleSeed           = uint64(1)
	quickLineCount       = 32
)

func main() {
	root := flag.String("root", ".", "path to repository root")
	flag.Parse()
	if err := run(*root); err != nil {
		log.Fatal(err)
	}
}

func run(root string) error {
	tokensSrc := filepath.Join(root, "third_party/razdel/razdel/tests/data/tokens.txt")
	sentsSrc := filepath.Join(root, "third_party/razdel/razdel/tests/data/sents.txt")
	outDir := filepath.Join(root, "testdata", "upstream")
	metaPath := filepath.Join(outDir, "META.txt")

	for _, p := range []string{tokensSrc, sentsSrc} {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("stat %s: %w", p, err)
		}
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	tf, err := os.Open(tokensSrc)
	if err != nil {
		return err
	}
	tokSample, err := fixture.SampleLines(tf, quickLineCount, sampleSeed)
	_ = tf.Close()
	if err != nil {
		return fmt.Errorf("sample tokens: %w", err)
	}
	if err := writeLines(filepath.Join(outDir, "quick_tokens_sample.txt"), tokSample); err != nil {
		return err
	}

	sf, err := os.Open(sentsSrc)
	if err != nil {
		return err
	}
	sentSample, err := fixture.SampleLines(sf, quickLineCount, sampleSeed)
	_ = sf.Close()
	if err != nil {
		return fmt.Errorf("sample sents: %w", err)
	}
	if err := writeLines(filepath.Join(outDir, "quick_sents_sample.txt"), sentSample); err != nil {
		return err
	}

	meta := fmt.Sprintf(`source_submodule_commit=%s
sample_algorithm=first_%d_indices_of_rand.Perm(n)_with_PCG(%d,%d)
quick_tokens_lines=%d
quick_sents_lines=%d
upstream_tokens_path=third_party/razdel/razdel/tests/data/tokens.txt
upstream_sents_path=third_party/razdel/razdel/tests/data/sents.txt
note=Python upstream uses random.seed(1) and random.sample(list(lines), count); Go uses PCG-based Perm for deterministic full-file sampling without loading all lines into a Python list first.
`, upstreamSubmoduleSHA, quickLineCount, sampleSeed, sampleSeed, len(tokSample), len(sentSample))

	if err := os.WriteFile(metaPath, []byte(meta), 0o644); err != nil {
		return err
	}
	log.Printf("wrote %s and samples (%d + %d lines)", outDir, len(tokSample), len(sentSample))
	return nil
}

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	for i, line := range lines {
		if i > 0 {
			if _, err := f.WriteString("\n"); err != nil {
				return err
			}
		}
		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}
	if len(lines) > 0 {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}
	return nil
}
