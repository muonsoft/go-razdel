# T012 — Sentenizer trivial rules

## Goal

Реализовать базовые trivial-правила sentenize перед переносом сложных исключений.

## Upstream references

- `third_party/razdel/razdel/segmenters/sentenize.py`:
  - `empty_side`
  - `no_space_prefix`
  - `lower_right`
  - `delimiter_right`

## Scope

### In scope
- Join, когда отсутствуют токены по краям.
- Join при отсутствии пробела справа.
- Join при `lowercase` правом токене.
- Join для delimiter-последователей и smile prefix.

### Out of scope
- Sokr/initials/list/bound/dash специфика.

## Steps

1. Перенести trivial функции в Go.
2. Встроить в порядок rule pipeline.
3. Добавить unit тесты на каждое правило и комбинации.

## Tests

- Unit: no-space scenarios.
- Unit: lowercase continuation.
- Unit: delimiter-right / smile cases.

## Acceptance criteria

- [x] Trivial слой полностью покрыт unit-тестами.
- [x] Порядок применения правил зафиксирован тестами.
- [x] `go test ./...` проходит.
