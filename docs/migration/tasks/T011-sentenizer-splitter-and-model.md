# T011 — Sentenizer splitter and split model

## Goal

Перенести основу sentenize: delimiter scanning, `SentSplit` и вычисляемые свойства.

## Upstream references

- `third_party/razdel/razdel/segmenters/sentenize.py`:
  - `SentSplit`
  - `SentSplitter`
  - regex: `DELIMITER`, `FIRST_TOKEN`, `LAST_TOKEN`, `WORD`

## Scope

### In scope
- Поиск delimiter (`.?!…`, quotes/brackets, smiles).
- Формирование `left/right` окна.
- Вычисляемые свойства (`left_token`, `right_token`, etc.).

### Out of scope
- Все rule-решения join/split.
- Post trim.

## Steps

1. Реализовать splitter и split-модель.
2. Проверить корректность токен-доступоров (`left_token`, `right_token`).
3. Добавить unit-тесты на структуру split событий.

## Tests

- Unit: пустые/пробельные строки.
- Unit: delimiter extraction.
- Unit: свойства split на строках с quotes/brackets/smiles.

## Acceptance criteria

- [ ] `SentSplitter` поведение соответствует upstream на unit-кейсах.
- [ ] Все ключевые cached-property эквиваленты покрыты тестами.
- [ ] `go test ./...` проходит.
