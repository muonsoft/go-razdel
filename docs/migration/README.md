# Go migration plan for `natasha/razdel`

Этот каталог содержит полный план миграции `third_party/razdel` на Go с упором на behavioral parity и проверяемые инкременты.

## Цели

1. Поведенчески совместимый `Tokenize` и `Sentenize`.
2. Контрактные гарантии API (`docs/contracts.md`) для span/offset.
3. Высокое тестовое покрытие на основе upstream unit/int тестов.
4. Контролируемое выполнение через маленькие независимые задачи.

## Структура

- `context/upstream-test-inventory.md` — аудит upstream-кода и тестов.
- `testing/strategy.md` — стратегия тестирования и quality gates.
- `execution/workflow.md` — порядок выполнения задач и инженерный цикл.
- `control/tracker.md` — трекер прогресса, зависимости и acceptance gates.
- `tasks/TXXX-*.md` — отдельный файл на каждую минимальную подзадачу.

## Глобальные quality gates

- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...`
- parity-тесты относительно upstream fixtures.

## Порядок реализации (высокоуровнево)

1. Подготовка тестовой инфраструктуры и фикстур parity.
2. Пошаговый перенос токенизации (от split до rule-by-rule).
3. Пошаговый перенос сегментации предложений (split + rule-by-rule + post).
4. Усиление тестового контура (corpus, differential, fuzz, benchmark).
5. Финальная стабилизация и фиксация известных расхождений.

## Task index

1. [T001 — Test harness bootstrap](tasks/T001-test-harness-bootstrap.md)
2. [T002 — Upstream fixtures ingestion](tasks/T002-upstream-fixtures-ingestion.md)
3. [T003 — Core API/offset invariants](tasks/T003-core-api-offset-invariants.md)
4. [T004 — Tokenizer atoms and splitter](tasks/T004-tokenizer-atoms-and-splitter.md)
5. [T005 — Tokenizer trivial split behavior](tasks/T005-tokenizer-trivial-split-behavior.md)
6. [T006 — Tokenizer dash/underscore joins](tasks/T006-tokenizer-dash-underscore-rules.md)
7. [T007 — Tokenizer float/fraction joins](tasks/T007-tokenizer-float-fraction-rules.md)
8. [T008 — Tokenizer punct/other/yahoo rules](tasks/T008-tokenizer-punct-other-yahoo-rules.md)
9. [T009 — Tokenizer unit parity suite](tasks/T009-tokenizer-unit-parity-suite.md)
10. [T010 — Tokenizer integration corpus suite](tasks/T010-tokenizer-integration-corpus-suite.md)
11. [T011 — Sentenizer splitter and split model](tasks/T011-sentenizer-splitter-and-model.md)
12. [T012 — Sentenizer trivial rules](tasks/T012-sentenizer-trivial-rules.md)
13. [T013 — Sentenizer sokr and initials rules](tasks/T013-sentenizer-sokr-initials-rules.md)
14. [T014 — Sentenizer quote/bracket bound rules](tasks/T014-sentenizer-quote-bracket-bound-rules.md)
15. [T015 — Sentenizer bullet and dash rules](tasks/T015-sentenizer-bullet-dash-rules.md)
16. [T016 — Sentenizer unit parity and post-processing API parity](tasks/T016-sentenizer-postprocessing-api-parity.md)
17. [T017 — Sentenizer integration corpus suite](tasks/T017-sentenizer-integration-corpus-suite.md)
18. [T018 — Differential/fuzz/benchmark hardening](tasks/T018-differential-fuzz-benchmark-hardening.md)
