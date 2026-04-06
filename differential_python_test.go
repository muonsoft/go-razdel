package razdel_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/fixture"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

// envRazdelDifferentialPython disables differential runs when set to "0".
const envRazdelDifferentialPython = "RAZDEL_DIFFERENTIAL_PYTHON"

type pyDiffCase struct {
	Mode string `json:"mode"`
	Text string `json:"text"`
}

type pyDiffBatchRequest struct {
	Cases []pyDiffCase `json:"cases"`
}

type pyDiffBatchResponse struct {
	Results [][]string `json:"results"`
}

var (
	pyDiffOnce   sync.Once
	pyDiffOK     bool
	pyDiffReason string
)

func pythonDifferentialReady(root string) (ok bool, reason string) {
	pyDiffOnce.Do(func() {
		if os.Getenv(envRazdelDifferentialPython) == "0" {
			pyDiffReason = envRazdelDifferentialPython + "=0"
			return
		}
		py, err := exec.LookPath("python3")
		if err != nil {
			pyDiffReason = "python3 not in PATH"
			return
		}
		razdelRoot := filepath.Join(root, "third_party", "razdel")
		if _, err := os.Stat(filepath.Join(razdelRoot, "razdel", "__init__.py")); err != nil {
			pyDiffReason = "third_party/razdel submodule not checked out"
			return
		}
		cmd := exec.Command(py, "-c", "import razdel")
		cmd.Env = append(os.Environ(), "PYTHONPATH="+razdelRoot)
		cmd.Dir = root
		if err := cmd.Run(); err != nil {
			pyDiffReason = "cannot import razdel with PYTHONPATH=" + razdelRoot + ": " + err.Error()
			return
		}
		pyDiffOK = true
	})
	return pyDiffOK, pyDiffReason
}

func runPythonDiffBatch(t *testing.T, root string, cases []pyDiffCase) [][]string {
	t.Helper()
	script := filepath.Join(root, "testdata", "python", "razdel_diff_runner.py")
	if _, err := os.Stat(script); err != nil {
		t.Fatalf("diff runner script: %v", err)
	}
	body, err := json.Marshal(pyDiffBatchRequest{Cases: cases})
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("python3", script)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PYTHONPATH="+filepath.Join(root, "third_party", "razdel"))
	cmd.Stdin = bytes.NewReader(body)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			t.Fatalf("python diff runner: %v\nstderr:\n%s", err, stderr.String())
		}
		t.Fatal(err)
	}
	var resp pyDiffBatchResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("decode python stdout: %v\nraw:\n%s", err, stdout.String())
	}
	if len(resp.Results) != len(cases) {
		t.Fatalf("python returned %d results, want %d", len(resp.Results), len(cases))
	}
	return resp.Results
}

func loadPartitionTexts(t *testing.T, root, rel string) []string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var texts []string
	for _, p := range fixture.ParsePartitionLines(string(raw)) {
		texts = append(texts, p.Text())
	}
	return texts
}

func goTokenTexts(s string) []string {
	toks := razdel.Tokenize(s)
	out := make([]string, 0, len(toks))
	for _, tok := range toks {
		out = append(out, tok.Text)
	}
	return out
}

func goSentenceTexts(s string) []string {
	sents := razdel.Sentenize(s)
	out := make([]string, 0, len(sents))
	for _, sent := range sents {
		out = append(out, sent.Text)
	}
	return out
}

func TestDifferential_tokenize_vsPython_quickSample(t *testing.T) {
	root := testkit.ModuleRoot(t)
	ok, why := pythonDifferentialReady(root)
	if !ok {
		t.Skip("differential vs Python:", why)
	}
	texts := loadPartitionTexts(t, root, "testdata/upstream/quick_tokens_sample.txt")
	cases := make([]pyDiffCase, len(texts))
	for i, s := range texts {
		cases[i] = pyDiffCase{Mode: "tokenize", Text: s}
	}
	wantPy := runPythonDiffBatch(t, root, cases)
	for i, s := range texts {
		got := goTokenTexts(s)
		if len(got) != len(wantPy[i]) {
			t.Fatalf("case %d text %q: len(go)=%d len(py)=%d\ngo: %#v\npy: %#v",
				i+1, s, len(got), len(wantPy[i]), got, wantPy[i])
		}
		for j := range got {
			if got[j] != wantPy[i][j] {
				t.Fatalf("case %d text %q token %d: go %q py %q", i+1, s, j, got[j], wantPy[i][j])
			}
		}
		testkit.AssertTokenOffsetContract(t, s, razdel.Tokenize(s))
	}
}

func TestDifferential_sentenize_vsPython_quickSample(t *testing.T) {
	root := testkit.ModuleRoot(t)
	ok, why := pythonDifferentialReady(root)
	if !ok {
		t.Skip("differential vs Python:", why)
	}
	texts := loadPartitionTexts(t, root, "testdata/upstream/quick_sents_sample.txt")
	cases := make([]pyDiffCase, len(texts))
	for i, s := range texts {
		cases[i] = pyDiffCase{Mode: "sentenize", Text: s}
	}
	wantPy := runPythonDiffBatch(t, root, cases)
	for i, s := range texts {
		got := goSentenceTexts(s)
		if len(got) != len(wantPy[i]) {
			t.Fatalf("case %d text %q: len(go)=%d len(py)=%d\ngo: %#v\npy: %#v",
				i+1, s, len(got), len(wantPy[i]), got, wantPy[i])
		}
		for j := range got {
			if got[j] != wantPy[i][j] {
				t.Fatalf("case %d text %q sentence %d: go %q py %q", i+1, s, j, got[j], wantPy[i][j])
			}
		}
		testkit.AssertSentenceOffsetContract(t, s, razdel.Sentenize(s))
	}
}
