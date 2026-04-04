# Testing strategy for Go port parity

## Основные принципы

1. **Parity first**: каждый переносимый кусок поведения связан с upstream-кейсом.
2. **Малые инкременты**: одна-две rule-группы за задачу.
3. **Контракт + поведение**: тесты одновременно проверяют offsets и разбиение.

## Тестовые слои

## L0 — Contract invariants (обязательно всегда)

- Проверить для каждого `Token`/`Sentence`:
  - `0 <= Start <= End <= len(text)`
  - `Text == text[Start:End]`
  - `Start`/`End` — byte offsets UTF-8.
- Никаких panic на корректном UTF-8.

## L1 — Unit parity by rule

- Табличные тесты по rule-группам:
  - Tokenize: split-space, dash/underscore, float/fraction, punct/other/yahoo.
  - Sentenize: trivial, sokr, bound, bullet, dash.
- Каждая строка testcase содержит upstream reference.

## L2 — Upstream unit fixtures (exact examples)

- Полная миграция `UNIT` наборов из:
  - `test_tokenize.py`
  - `test_sentenize.py`
- Тесты должны быть детерминированными и быстрыми.

## L3 — Integration corpus parity

- Runner по форматам `tokens.txt` / `sents.txt`:
  - `quick`: deterministic sample (фиксированный seed + count).
  - `full`: весь файл (опционально, не по умолчанию).
- PR-гейт: quick.
- Полный прогон: nightly/manual.

## L4 — Differential tests (Go vs Python upstream)

- Для sampled кейсов запускать Python `razdel` из submodule и сравнивать результат.
- Тесты делать `t.Skip(...)`, если Python окружение недоступно.

## L5 — Robustness/quality

- Fuzz:
  - `FuzzTokenizeNoPanic`
  - `FuzzSentenizeNoPanic`
  - проверка span-инвариантов.
- Benchmarks:
  - короткий текст;
  - длинный текст;
  - Unicode-heavy;
  - punctuation-heavy.

## Организация testdata

Рекомендуемая структура:

- `testdata/upstream/tokenize/unit.txt`
- `testdata/upstream/sentenize/unit.txt`
- `testdata/upstream/tokenize/tokens.sample.txt`
- `testdata/upstream/sentenize/sents.sample.txt`
- `testdata/upstream/*.meta.json` (seed, count, commit SHA submodule)

## Команды проверки

- Минимум на каждую итерацию:
  - `go test ./...`
  - `go vet ./...`
  - `golangci-lint run ./...`
- Дополнительно (по флагам/тегам):
  - full corpus;
  - differential;
  - benchmarks.

## Критерий “coverage high-quality”

Считаем покрытие достаточным, когда:

1. Все upstream unit кейсы перенесены и проходят.
2. Quick corpus parity стабилен в CI.
3. Full corpus parity не показывает систематических расхождений.
4. Fuzz не выявляет panic/нарушений инвариантов.
5. Известные отклонения задокументированы с примерами.
