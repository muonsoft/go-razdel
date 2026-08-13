# Contributing

Thanks for looking into `go-razdel`. Please open an [issue](https://github.com/muonsoft/go-razdel/issues) before large changes, especially anything that would change token/sentence **texts** relative to Python `razdel`.

## Development setup

```bash
git clone https://github.com/muonsoft/go-razdel.git
cd go-razdel
git submodule update --init --recursive
```

The submodule `third_party/razdel` is the behavioral etalon. Do not edit it in a feature PR; submodule bumps are a separate change with old and new commit SHAs.

Minimum Go version is in `go.mod` (currently 1.24).

## Checks

```bash
go test ./...
go vet ./...
golangci-lint run ./...
```

CI (`.github/workflows/ci.yml`) runs the same set and **requires** the Python differential (`RAZDEL_DIFFERENTIAL_PYTHON=1`). Locally the differential is skipped if `python3` cannot import pinned `razdel`.

### Python differential

Needs `python3` on `PATH` and an initialized submodule. Tests set `PYTHONPATH` themselves.

```bash
# force fail (same as CI) if Python razdel is missing
RAZDEL_DIFFERENTIAL_PYTHON=1 go test ./...

# skip
RAZDEL_DIFFERENTIAL_PYTHON=0 go test ./...
```

### Full corpora

Quick samples run with `go test ./...`. The full upstream files (~209k tokenize / ~49k sentenize lines) are manual or the nightly workflow `.github/workflows/corpus.yml`:

```bash
RAZDEL_TOKENIZE_INTEGRATION_FULL=1 go test ./internal/tokenize/ -run TestIntegration_tokenize_full_corpus -count=1
RAZDEL_SENTENIZE_INTEGRATION_FULL=1 go test ./internal/sentenize/ -run TestIntegration_sentenize_full_corpus -count=1
```

Fixture format and drift notes: [`testdata/upstream/README.md`](testdata/upstream/README.md).

### Fuzz and benchmarks

```bash
go test -fuzz=FuzzTokenize  -fuzztime=30s .
go test -fuzz=FuzzSentenize -fuzztime=30s .
go test -bench 'BenchmarkTokenize|BenchmarkSentenize' -benchmem .
```

Performance policy and baseline numbers: [`docs/performance.md`](docs/performance.md). PRs that touch hot paths get a `benchstat` comment from `.github/workflows/benchstat.yml`.

## Parity rules

1. Prefer a small, testable slice of behavior over a large rewrite.
2. Point at an upstream test or example (`third_party/razdel`) for the expected texts.
3. Compare Go vs Python by token/sentence **text**, not by numeric offsets (bytes vs code points).
4. Document new deviations in `README.md` and `docs/contracts.md` with an input example.
5. Public API and byte-offset contract: [`docs/contracts.md`](docs/contracts.md).

The original migration task list is complete; it remains as history under [`docs/migration/`](docs/migration/README.md).

## Releases

1. Put user-visible notes in the `## [Unreleased]` section of [`CHANGELOG.md`](CHANGELOG.md).
2. From `main`, run **Actions → Release → Run workflow** with a version like `v0.1.0` (or push a `v*.*.*` tag).
3. The workflow tests, moves `[Unreleased]` into `[x.y.z]`, commits the changelog on `main`, and publishes a GitHub Release.

If the version section already exists in `CHANGELOG.md`, the workflow only publishes the release.

## User-facing docs

Keep [`README.md`](README.md) (English) and [`README.ru.md`](README.ru.md) in sync when examples, the public API, or compatibility notes change.

## Repository settings (maintainers)

About, topics, and related GitHub UI fields are listed in [`docs/github.md`](docs/github.md). Issue labels are defined in [`.github/labels.yml`](.github/labels.yml) and synced by `.github/workflows/label-sync.yml`.
