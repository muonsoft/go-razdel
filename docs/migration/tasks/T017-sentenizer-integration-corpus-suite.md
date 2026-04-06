# T017 — Sentenizer integration corpus suite

## Goal

Включить интеграционные parity-тесты sentenize на большом upstream корпусе.

## Upstream references

- `third_party/razdel/razdel/tests/test_sentenize.py` (`test_int`, `int_tests`)
- `third_party/razdel/razdel/tests/data/sents.txt` (48,735 строк)

## Scope

### In scope
- Quick deterministic sample режим для PR.
- Full corpus режим для manual/nightly.
- Диагностика несовпадений (diff format).

### Out of scope
- Изменения upstream datasets.

## Steps

1. Реализовать corpus runner для partition строк sentenize.
2. Добавить флаги quick/full и отчеты.
3. Подключить quick режим в CI.

## Tests

- Integration quick: sample parity.
- Integration full: manual/nightly validation.
- Regression: reproducible sampling.

## Acceptance criteria

- [x] Quick parity стабилен и быстрый.
- [x] Full прогон доступен и документирован.
- [x] Диффы читаемы и помогают triage расхождений.
- [x] `go test ./...` проходит в quick режиме.
