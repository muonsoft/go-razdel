# T004 — Tokenizer atoms and splitter

## Goal

Перенести базовую атомизацию и split-модель токенизатора (без всех join-правил).

## Upstream references

- `third_party/razdel/razdel/segmenters/tokenize.py`:
  - `ATOM`
  - `Atom`
  - `TokenSplit`
  - `TokenSplitter`

## Scope

### In scope
- Эквивалент классификации атомов (`RU/LAT/INT/PUNCT/OTHER`).
- Формирование split-окон (`left_n`/`right_n` эквиваленты).
- Корректная обработка delimiter между атомами.

### Out of scope
- Правила join/split поведения.

## Steps

1. Реализовать атомайзер по regex/классам символов.
2. Реализовать splitter с контекстным окном.
3. Добавить минимальные тесты на структуру split и типы атомов.

## Tests

- Unit: классификация `mβж`, `Δσ`, `1`, `,`, `...`.
- Unit: формирование split delimiter на пробелах/без пробелов.
- Contract: корректные исходные span/text.

## Acceptance criteria

- [ ] Atom типы и splitter поведение совпадают с upstream на unit-наборе.
- [ ] Контекстные left/right окна тестируются.
- [ ] Тесты независимы от финальной реализации join-правил.
- [ ] `go test ./...` проходит.
