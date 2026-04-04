# T006 — Tokenizer dash and underscore rules

## Goal

Перенести правила join для дефиса/тире и подчеркивания (Rule2112 ветка).

## Upstream references

- `third_party/razdel/razdel/segmenters/tokenize.py`:
  - `Rule2112`
  - `DashRule`
  - `UnderscoreRule`
- `test_tokenize.py` unit cases:
  - `что-то`
  - `К_тому_же`

## Scope

### In scope
- Join через `DASHES` при непунктуационных соседях.
- Join через `_` при непунктуационных соседях.

### Out of scope
- Float/fraction.
- Punct/other/yahoo исключения.

## Steps

1. Реализовать Rule2112-контур в Go.
2. Добавить DashRule + UnderscoreRule.
3. Проверить приоритет и порядок правил.

## Tests

- Unit: `что-то` как один токен.
- Unit: `К_тому_же` как один токен.
- Negative: дефис между пунктуацией не склеивается.

## Acceptance criteria

- [ ] Upstream unit примеры для dash/underscore проходят.
- [ ] Есть negative tests на недопустимый join.
- [ ] `go test ./...` проходит.
