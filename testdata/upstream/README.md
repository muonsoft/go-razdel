# Upstream-derived test fixtures (T002)

These files mirror the `UNIT` lists from upstream:

- `unit_tokenize.txt` — `third_party/razdel/razdel/tests/test_tokenize.py`
- `unit_sentenize.txt` — `third_party/razdel/razdel/tests/test_sentenize.py`

Each non-empty line is one partition record: chunks separated by ASCII `|`, same as `parse_partition` in upstream `partition.py`. Whitespace-only chunks are **fill** gaps between expected segments.

The upstream tokenizer includes an empty-string case (`''`). Because fixture files are newline-delimited, that case is encoded as a single line `#empty` (see `internal/fixture.EmptyPartitionMarker`).

## Quick corpus samples

- `quick_tokens_sample.txt` / `quick_sents_sample.txt` — deterministic subsets of `razdel/tests/data/tokens.txt` and `sents.txt`.
- `META.txt` — pinned submodule commit SHA, seed, line counts, and sampling notes.

Regenerate quick samples after updating the submodule (or changing counts/seed):

```bash
go run ./tools/genupstreamfixtures
```

Upstream integration tests use `random.seed(1)` and `random.sample`. The Go port uses `math/rand/v2` with a fixed PCG seed and the first *k* indices of `Perm(n)` so the entire corpus can be sampled without materializing a Python-sized list in memory.
