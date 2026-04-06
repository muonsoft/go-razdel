package fixture

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/muonsoft/go-razdel/internal/testkit"
)

func TestIsFill(t *testing.T) {
	t.Parallel()
	if !IsFill(" ") || !IsFill("  ") || !IsFill("\t") {
		t.Fatal("whitespace-only chunks should be fill")
	}
	if IsFill("a") || IsFill(" a") {
		t.Fatal("non-fill chunks")
	}
}

func TestParsePartitionLine_emptyStringUpstream(t *testing.T) {
	t.Parallel()
	p := ParsePartitionLine("")
	if len(p.Chunks) != 1 || p.Chunks[0] != "" {
		t.Fatalf("want single empty chunk like Python '', got %#v", p.Chunks)
	}
	segs := p.ExpectedSegments()
	if len(segs) != 0 {
		t.Fatalf("empty partition: want no segments, got %d", len(segs))
	}
}

func TestParsePartitionLines_emptyMarker(t *testing.T) {
	t.Parallel()
	parts := ParsePartitionLines("#empty\n")
	if len(parts) != 1 {
		t.Fatalf("want 1 partition, got %d", len(parts))
	}
	if len(parts[0].Chunks) != 1 || parts[0].Chunks[0] != "" {
		t.Fatalf("marker should yield [''], got %#v", parts[0].Chunks)
	}
}

func TestExpectedSegments_tokenizeUnitFile(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "testdata", "upstream", "unit_tokenize.txt")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parts := ParsePartitionLines(string(raw))
	if len(parts) != 14 {
		t.Fatalf("unit_tokenize: want 14 partitions, got %d", len(parts))
	}
	// Last line is #empty
	last := parts[len(parts)-1]
	if len(last.ExpectedSegments()) != 0 {
		t.Fatalf("empty case: want 0 segments")
	}
	// mβж — two runes, one segment
	mbeta := parts[11]
	segs := mbeta.ExpectedSegments()
	if len(segs) != 1 || segs[0].Text != "mβж" || segs[0].StartRune != 0 || segs[0].EndRune != 3 {
		t.Fatalf("mβж segment: got %+v", segs)
	}
	// »||. — three chunks, middle is fill
	gt := parts[6]
	segs = gt.ExpectedSegments()
	wantText := []string{"»", "."}
	if len(segs) != len(wantText) {
		t.Fatalf("»||.: want %d segments, got %+v", len(wantText), segs)
	}
	for i := range wantText {
		if segs[i].Text != wantText[i] {
			t.Fatalf("segment %d: want %q, got %q", i, wantText[i], segs[i].Text)
		}
	}
}

func TestExpectedSegments_sentenizeUnitFile(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "testdata", "upstream", "unit_sentenize.txt")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parts := ParsePartitionLines(string(raw))
	if len(parts) != 34 {
		t.Fatalf("unit_sentenize: want 34 partitions, got %d", len(parts))
	}
	// Line with "| |" gap between two sentences
	line := parts[1]
	segs := line.ExpectedSegments()
	if len(segs) != 2 {
		t.Fatalf("want 2 sentences, got %+v", segs)
	}
	if segs[0].Text == "" || segs[1].Text == "" {
		t.Fatal("non-empty sentence texts expected")
	}
	whole := line.Text()
	if segs[0].StartRune != 0 {
		t.Fatalf("first start: got %d", segs[0].StartRune)
	}
	runes := []rune(whole)
	got0 := string(runes[segs[0].StartRune:segs[0].EndRune])
	got1 := string(runes[segs[1].StartRune:segs[1].EndRune])
	if got0 != segs[0].Text || got1 != segs[1].Text {
		t.Fatalf("rune slice mismatch: %q %q vs %q %q", got0, got1, segs[0].Text, segs[1].Text)
	}
}

func TestSampleLines_reproducible(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	path := filepath.Join(root, "third_party", "razdel", "razdel", "tests", "data", "tokens.txt")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	const count = 32
	const seed = uint64(1)
	a, err := SampleLines(f, count, seed)
	if err != nil {
		t.Fatal(err)
	}
	f2, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f2.Close() }()
	b, err := SampleLines(f2, count, seed)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("same seed should yield identical sample")
	}
}

func TestQuickSamples_matchGenerated(t *testing.T) {
	t.Parallel()
	root := testkit.ModuleRoot(t)
	tokensPath := filepath.Join(root, "testdata", "upstream", "quick_tokens_sample.txt")
	sentsPath := filepath.Join(root, "testdata", "upstream", "quick_sents_sample.txt")
	wantTok, err := os.ReadFile(tokensPath)
	if err != nil {
		t.Fatal(err)
	}
	wantSent, err := os.ReadFile(sentsPath)
	if err != nil {
		t.Fatal(err)
	}
	srcTok, err := os.Open(filepath.Join(root, "third_party", "razdel", "razdel", "tests", "data", "tokens.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = srcTok.Close() }()
	gotTokLines, err := SampleLines(srcTok, 32, 1)
	if err != nil {
		t.Fatal(err)
	}
	gotTok := strings.Join(append([]string{}, gotTokLines...), "\n") + "\n"
	if gotTok != string(wantTok) {
		t.Fatal("quick_tokens_sample.txt out of sync with SampleLines; run go run ./tools/genupstreamfixtures")
	}
	srcSent, err := os.Open(filepath.Join(root, "third_party", "razdel", "razdel", "tests", "data", "sents.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = srcSent.Close() }()
	gotSentLines, err := SampleLines(srcSent, 32, 1)
	if err != nil {
		t.Fatal(err)
	}
	gotSent := strings.Join(append([]string{}, gotSentLines...), "\n") + "\n"
	if gotSent != string(wantSent) {
		t.Fatal("quick_sents_sample.txt out of sync; run go run ./tools/genupstreamfixtures")
	}
}
