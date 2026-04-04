# T018 — Differential, fuzz, and benchmark hardening

## Goal

Закрепить качество через differential сравнение с Python upstream, fuzz и benchmark набор.

## Upstream references

- `third_party/razdel` (Python runtime behavior)
- `docs/contracts.md` (инварианты)
- `third_party/razdel/razdel/tests/data/*.txt` (источник входов)

## Scope

### In scope
- Differential tests: Go output vs Python output на sampled input.
- Fuzz tests: отсутствие panic + offsets invariants.
- Benchmarks: короткие/длинные/Unicode-heavy/punct-heavy тексты.

### Out of scope
- Оптимизация производительности вне выявленных bottleneck.

## Steps

1. Добавить go tests, вызывающие upstream Python (с graceful skip).
2. Добавить fuzz targets для `Tokenize` и `Sentenize`.
3. Добавить benchmark suite и baseline результаты.
4. Задокументировать процедуру запуска в CI/manual.

## Tests

- Differential: sample parity tokenize/sentenize.
- Fuzz: инварианты на случайном вводе.
- Benchmark: `go test -bench .` для ключевых сценариев.

## Acceptance criteria

- [ ] Differential тесты доступны и воспроизводимы.
- [ ] Fuzz не выявляет panic/контрактных нарушений.
- [ ] Benchmarks задокументированы и повторяемы.
- [ ] `go test ./...` проходит (без обязательного Python в окружении).
