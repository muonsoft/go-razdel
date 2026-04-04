# T016 — Sentenizer unit parity and post-processing API parity

## Goal

Закрыть полный unit parity для sentenize и довести API-паритет: корректный post-step (`strip`) и span/text соответствие.

## Upstream references

- `third_party/razdel/razdel/segmenters/sentenize.py`:
  - `SentSegmenter.post` (`yield chunk.strip()`)
- `third_party/razdel/razdel/tests/test_sentenize.py` (`UNIT`)
- `docs/contracts.md`

## Scope

### In scope
- Полный перенос upstream `UNIT` набора для sentenize.
- Post processing после сегментации.
- Корректное восстановление offsets после trimming.
- Финальный API уровень для `Sentenize`.

### Out of scope
- Интеграционные full-corpus тесты (это T017).

## Steps

1. Перенести все `UNIT` кейсы из upstream sentenize в Go fixtures/table tests.
2. Внедрить `post`-обработку в pipeline.
3. Проверить, что offsets остаются валидными после trim.
4. Добавить API-level regression тесты.

## Tests

- Unit parity: все кейсы `UNIT` из `test_sentenize.py`.
- Unit: leading/trailing spaces around sentence chunks.
- Contract: `Text == text[Start:End]` после post-step.
- Regression: пустые/пробельные входы.

## Acceptance criteria

- [ ] 100% upstream sentenize `UNIT` кейсов присутствуют и проходят.
- [ ] Поведение `strip()` parity подтверждено.
- [ ] API контракт `Sentence` соблюдается на всех тестах.
- [ ] `go test ./...` проходит.
