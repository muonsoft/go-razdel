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

## Tokenizer integration corpus (T010)

- **Quick (default):** `go test ./...` runs `TestIntegration_tokenize_quick_corpus` against `quick_tokens_sample.txt` (deterministic sample from upstream `tokens.txt`).
- **Full (manual / nightly):** run the entire upstream file `third_party/razdel/razdel/tests/data/tokens.txt` (~209k lines):

```bash
RAZDEL_TOKENIZE_INTEGRATION_FULL=1 go test ./internal/tokenize/ -run TestIntegration_tokenize_full_corpus -count=1 -v
```

Set `RAZDEL_TOKENIZE_INTEGRATION_FULL=1` to enable; any other value (or unset) skips the full test.

The full run compares `tokenize.TokenTexts` to the **partition etalon** in `tokens.txt` (same contract as upstream `test_int`). **25** lines in the pinned submodule disagree with `razdel.segmenters.tokenize` on the same text (fixture drift). Those partition strings are listed in `internal/tokenize/upstream_tokens_txt_drift_test.go`; the full test **skips** them (records `t.Skip` subtests with a reason in `-v` output) and asserts the rest of the corpus. After updating the submodule, refresh that list if the drift set changes.

## Sentenizer integration corpus (T017)

- **Quick (default):** `go test ./...` runs `TestIntegration_sentenize_quick_corpus` against `quick_sents_sample.txt` (deterministic sample from upstream `sents.txt`).
- **Full (manual / nightly):** run the entire upstream file `third_party/razdel/razdel/tests/data/sents.txt` (~48.7k lines):

```bash
RAZDEL_SENTENIZE_INTEGRATION_FULL=1 go test ./internal/sentenize/ -run TestIntegration_sentenize_full_corpus -count=1 -v
```

Set `RAZDEL_SENTENIZE_INTEGRATION_FULL=1` to enable; any other value (or unset) skips the full test.

The full run compares sentence texts from `razdel.Sentenize` to the **partition etalon** in `sents.txt` (same contract as upstream `test_sentenize.test_int`). If the pinned submodule ever contains lines where the file disagrees with `razdel.segmenters.sentenize` on the same text, list those partition strings in `internal/sentenize/upstream_sents_txt_drift_test.go` (same pattern as tokenizer drift). With the current pin, the drift set is empty.
