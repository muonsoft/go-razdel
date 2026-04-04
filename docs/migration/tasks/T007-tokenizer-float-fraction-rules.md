# T007 — Tokenizer float and fraction rules

## Goal

Реализовать join для числовых паттернов `INT delimiter INT`.

## Upstream references

- `third_party/razdel/razdel/segmenters/tokenize.py`:
  - `FloatRule`
  - `FractionRule`
- `test_tokenize.py` unit:
  - `1,5`
  - `1/2`

## Scope

### In scope
- Join по `.` и `,` между двумя INT.
- Join по `/` и `\` между двумя INT.

### Out of scope
- Общая нормализация чисел.
- Локалезависимые форматы помимо upstream логики.

## Steps

1. Добавить FloatRule.
2. Добавить FractionRule.
3. Протестировать конфликтные кейсы рядом с пунктуацией.

## Tests

- Unit: `1,5` / `1.5` / `1/2` / `1\2`.
- Negative: `a,b` не должен склеиваться как число.
- Regression: соседство с другими правилами не ломает dash/underscore.

## Acceptance criteria

- [ ] Numeric join поведение совпадает с upstream.
- [ ] Есть позитивные и негативные кейсы.
- [ ] `go test ./...` проходит.
