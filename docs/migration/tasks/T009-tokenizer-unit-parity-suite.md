# T009 — Tokenizer unit parity suite

## Goal

Собрать полный unit parity suite для токенизации по upstream `UNIT`.

## Upstream references

- `third_party/razdel/razdel/tests/test_tokenize.py` (`UNIT`)
- `third_party/razdel/razdel/tests/partition.py`

## Scope

### In scope
- Все примеры `UNIT` из upstream как table-driven Go tests.
- Единый тест-раннер `partition -> expected tokens`.
- Подробный diff в ошибках (как в `pytest_assertrepr_compare`).

### Out of scope
- Большой интеграционный корпус `tokens.txt`.

## Steps

1. Перенести все `UNIT` строки в go-fixture.
2. Настроить table tests с информативным именованием кейсов.
3. Добавить удобный output diff для отладки расхождений.

## Tests

- Unit parity: все кейсы `UNIT`.
- Contract: offsets валидны в каждом кейсе.
- Stability: deterministic order/результат.

## Acceptance criteria

- [ ] 100% upstream UNIT tokenize кейсов присутствуют.
- [ ] Все кейсы проходят.
- [ ] Сообщения о падениях позволяют быстро локализовать rule.
- [ ] `go test ./...` проходит.
