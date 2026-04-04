# T014 — Sentenizer quote and bracket bound rules

## Goal

Перенести логику границ предложений вокруг кавычек и закрывающих скобок.

## Upstream references

- `third_party/razdel/razdel/segmenters/sentenize.py`:
  - `close_bound`
  - `close_quote`
  - `close_bracket`
- `third_party/razdel/razdel/segmenters/punct.py`
- `test_sentenize.py` quote-bound кейсы.

## Scope

### In scope
- Различение open/close/generic quotes.
- Join/split decisions для закрывающих скобок.
- Проверка зависимостей от левого ending token.

### Out of scope
- Sokr и list item логика.

## Steps

1. Перенести наборы символов quotes/brackets.
2. Реализовать close-bound правила.
3. Добавить unit и regression кейсы на пограничные цитаты.

## Tests

- Unit: `"...".` и похожие конструкции.
- Unit: delimiter в generic quotes (`"`/`'`/`„`).
- Regression: не ломать базовые trivial rules.

## Acceptance criteria

- [ ] Boundary в кавычках/скобках совпадает с upstream.
- [ ] Есть coverage для разных типов кавычек.
- [ ] `go test ./...` проходит.
