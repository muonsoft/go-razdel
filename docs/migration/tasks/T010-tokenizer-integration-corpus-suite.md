# T010 — Tokenizer integration corpus suite

## Goal

Подключить интеграционное parity-тестирование токенизации на большом upstream корпусе.

## Upstream references

- `third_party/razdel/razdel/tests/test_tokenize.py` (`test_int`, `int_tests`)
- `third_party/razdel/razdel/tests/data/tokens.txt` (208,995 строк)

## Scope

### In scope
- Quick mode: deterministic sample из `tokens.txt`.
- Full mode: опциональный прогон всего `tokens.txt`.
- Метрики и отчет по числу несовпадений.

### Out of scope
- Постоянный full run на каждом PR.

## Steps

1. Сделать test runner для partition-строк корпуса.
2. Добавить флаги/переменные окружения для quick/full режимов.
3. Подготовить CI-конфигурацию для quick режима.

## Tests

- Integration quick: sample parity.
- Integration full (manual/nightly): 100% корпуса.
- Regression: стабильность sampling.

## Acceptance criteria

- [ ] Quick режим включен и стабилен.
- [ ] Full режим доступен и документирован.
- [ ] Отчеты о расхождениях воспроизводимы.
- [ ] `go test ./...` проходит в quick режиме.
