# T002 — Upstream fixtures ingestion

## Goal

Сделать воспроизводимый импорт и использование upstream фикстур для тестов Go.

## Upstream references

- `third_party/razdel/razdel/tests/test_tokenize.py`
- `third_party/razdel/razdel/tests/test_sentenize.py`
- `third_party/razdel/razdel/tests/data/tokens.txt`
- `third_party/razdel/razdel/tests/data/sents.txt`

## Scope

### In scope
- Перенос UNIT-кейсов в локальные go-friendly fixture файлы.
- Скрипт/утилита deterministic sampling из больших датасетов (seed фиксирован).
- Метаданные источника (submodule SHA, seed, count).

### Out of scope
- Изменение содержимого upstream submodule.
- Полный прогон корпуса в каждый CI job.

## Steps

1. Сформировать `testdata/upstream/...` с unit fixtures.
2. Подготовить sampled fixtures для quick CI режимов.
3. Добавить README/meta с происхождением фикстур и правилами обновления.

## Tests

- Unit: parser fixture-файлов.
- Integration: quick sample читается и прогоняется тест-раннером.
- Reproducibility: один и тот же seed дает одинаковый sample.

## Acceptance criteria

- [ ] Unit fixtures покрывают upstream UNIT кейсы.
- [ ] Quick sample fixtures детерминированы.
- [ ] Источник и параметры выборки задокументированы.
- [ ] `go test ./...` проходит.
