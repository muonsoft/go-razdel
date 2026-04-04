# T003 — Core API and offset invariants

## Goal

Закрепить контракт API и byte-offset семантику тестами до переноса алгоритмов.

## Upstream references

- `docs/contracts.md`
- `third_party/razdel/razdel/substring.py` (`find_substrings` логика позиционирования)

## Scope

### In scope
- Инварианты `Start/End/Text` для `Token`/`Sentence`.
- Тесты на UTF-8 многобайтовые символы.
- Детеминированность и отсутствие side effects.

### Out of scope
- Точное правило-правило parity tokenize/sentenize.

## Steps

1. Добавить table tests на offsets для нескольких классов входов.
2. Добавить helper проверки half-open интервала `[Start, End)`.
3. Зафиксировать ожидаемое поведение пустого ввода.

## Tests

- Contract: инварианты для ASCII/кириллицы/emoji.
- Regression: пустой ввод возвращает пустой результат.
- Stability: повторный вызов дает идентичный результат.

## Acceptance criteria

- [x] Контракт offsets строго зафиксирован тестами.
- [x] UTF-8 кейсы присутствуют.
- [x] Нет panic на корректном UTF-8.
- [x] `go test ./...` проходит.
