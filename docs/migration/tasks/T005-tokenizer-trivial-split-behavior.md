# T005 — Tokenizer trivial split behavior

## Goal

Реализовать базовое поведение split по delimiter-пробелам (trivial rule).

## Upstream references

- `third_party/razdel/razdel/segmenters/tokenize.py`:
  - `split_space`
  - `TokenSegmenter.segment`

## Scope

### In scope
- Базовое разделение токенов по delimiter при наличии пробела.
- Базовый цикл сегментации (buffer + next part).

### Out of scope
- Join-исключения для дефисов/чисел/пунктуации.

## Steps

1. Внедрить trivial правило split-space.
2. Добавить базовую сегментацию без специальных join.
3. Покрыть кейсы: простые слова, пунктуация с пробелом.

## Tests

- Unit: `"a b"` => `["a","b"]`.
- Unit: `"привет, мир"` без спец join (ожидаемо с разделением по пробелу).
- Contract: spans корректны.

## Acceptance criteria

- [ ] Trivial split работает как минимальная база для последующих правил.
- [ ] Нет регрессий в contract tests.
- [ ] `go test ./...` проходит.
