# T013 — Sentenizer sokr and initials rules

## Goal

Перенести правила сокращений и инициалов для корректных границ предложений.

## Upstream references

- `third_party/razdel/razdel/segmenters/sentenize.py`:
  - `sokr_left`
  - `inside_pair_sokr`
  - `initials_left`
- `third_party/razdel/razdel/segmenters/sokr.py`
- `test_sentenize.py` cases (`т. д.`, `к.п.н.`, `т.е.`, `т.н.`)

## Scope

### In scope
- Импорт/репрезентация SOKR словарей.
- Pair сокращения и initials.
- Поведение с точками в сокращениях.

### Out of scope
- Quote/bracket/list/dash boundary rules.

## Steps

1. Перенести словари `SOKRS`, `PAIR_SOKRS`, `INITIALS`.
2. Реализовать 3 правила и включить в pipeline.
3. Подготовить targeted regression tests.

## Tests

- Unit: отдельные кейсы на каждое правило.
- Unit parity: соответствующие примеры из upstream `UNIT`.
- Regression: исключить ложные join после обычных слов.

## Acceptance criteria

- [ ] Сокращения и инициалы ведут себя как в upstream.
- [ ] Словари имеют тесты целостности/наличия ключевых элементов.
- [ ] `go test ./...` проходит.
