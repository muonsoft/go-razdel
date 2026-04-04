package fixture

import (
	"bufio"
	"io"
	"math/rand/v2"
)

// SampleLines reads all non-empty lines from r, then returns the first count lines of a
// deterministic permutation: rng.Perm(n)[:count] with PCG(seed, seed). Regenerate quick fixtures
// with tools/genupstreamfixtures.
func SampleLines(r io.Reader, count int, seed uint64) ([]string, error) {
	var lines []string
	sc := bufio.NewScanner(r)
	// Corpus lines can be long; default buf may be too small for sents.txt.
	const maxScan = 12 * 1024 * 1024
	buf := make([]byte, maxScan)
	sc.Buffer(buf, maxScan)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	n := len(lines)
	if count > n {
		count = n
	}
	if count == 0 {
		return nil, nil
	}
	rng := rand.New(rand.NewPCG(seed, seed))
	perm := rng.Perm(n)
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, lines[perm[i]])
	}
	return out, nil
}
