# T008 — Tokenizer punct/other/yahoo rules

## Goal

Перенести исключающие join-правила для пунктуации, mixed scripts и `yahoo!`.

## Upstream references

- `third_party/razdel/razdel/segmenters/tokenize.py`:
  - `punct`
  - `other`
  - `yahoo`
  - `SMILE`
- `test_tokenize.py` unit:
  - `...`
  - `:)))`
  - `mβж`
  - `Δσ`

## Scope

### In scope
- Склейка серий пунктуации (`...`, `?!`, `--`, `***`).
- Склейка smile-последовательностей.
- Склейка mixed-script OTHER/RU/LAT пар.
- Спец-кейс `yahoo!`.

### Out of scope
- Любые неописанные upstream эвристики.

## Steps

1. Реализовать `punct` rule с учетом smile/endings.
2. Реализовать `other` rule для Unicode-миксов.
3. Реализовать `yahoo` exception.

## Tests

- Unit: перенос кейсов из upstream UNIT.
- Regression: порядок правил не ломает numeric/dash поведение.
- Contract: spans корректны при сложной пунктуации.

## Acceptance criteria

- [ ] Все пунктуационные и mixed-script unit кейсы parity проходят.
- [ ] Набор правил работает в правильном порядке.
- [ ] `go test ./...` проходит.
