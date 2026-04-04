# T001 — Test harness bootstrap

## Goal

Подготовить базовую инфраструктуру Go-тестов для parity-переноса без переноса алгоритмов.

## Upstream references

- `third_party/razdel/razdel/tests/common.py`
- `third_party/razdel/razdel/tests/partition.py`
- `third_party/razdel/razdel/tests/conftest.py`

## Scope

### In scope
- Общие test helpers для сравнения expected vs actual сегментов.
- Парсер partition-формата с разделителем `|`.
- Единый assertion helper для span/text инвариантов.

### Out of scope
- Реализация tokenizer/sentenizer алгоритмов.
- Полный импорт больших датасетов.

## Steps

1. Создать `internal/testkit` (или эквивалент) с parser/assert helpers.
2. Добавить helper проверки контракта offsets для любых `[]Token`/`[]Sentence`.
3. Подключить helpers в существующие тесты-заглушки.

## Tests

- Unit: парсинг строки partition в ожидаемые сегменты.
- Contract: проверка byte offsets на ASCII и UTF-8 строках.
- Smoke: helpers используются из хотя бы одного теста.

## Acceptance criteria

- [ ] Partition parser работает для кейсов с `|`.
- [ ] Есть повторно используемые assert helper-ы для parity.
- [ ] Контракт offsets покрыт отдельным тестом.
- [ ] `go test ./...` проходит.
