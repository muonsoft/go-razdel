# Migration control tracker

Обновляйте статус по мере выполнения. Разрешен только один активный `In Progress` task.

## Legend

- `Planned` — готово к старту.
- `In Progress` — в работе.
- `Done` — выполнено и проверено.
- `Blocked` — есть внешний блокер.

## Task board

| ID | Task | Depends on | Status | Acceptance checkpoint |
|---|---|---|---|---|
| T001 | Test harness bootstrap | - | Done | Базовая test-инфраструктура создана |
| T002 | Upstream fixtures ingestion | T001 | Done | Unit/corpus фикстуры доступны в Go |
| T003 | Core API/offset invariants | T001 | Done | Инварианты offsets закреплены тестами |
| T004 | Tokenizer atoms and splitter | T003 | Done | Atom/split behavior совпадает с upstream |
| T005 | Tokenizer trivial split behavior | T004 | Done | split-space parity подтвержден |
| T006 | Tokenizer dash/underscore rules | T005 | Done | Кейсы `что-то`, `К_тому_же` проходят |
| T007 | Tokenizer float/fraction rules | T005 | Done | Кейсы `1,5`, `1/2` проходят |
| T008 | Tokenizer punct/other/yahoo rules | T006,T007 | Done | Пунктуация/smile/mixed-script parity |
| T009 | Tokenizer unit parity suite | T008 | Done | Все upstream unit tokenize кейсы в Go |
| T010 | Tokenizer integration corpus suite | T009,T002 | Done | quick/full corpus режимы для tokenize |
| T011 | Sentenizer splitter and model | T003 | Done | SentSplit/SentSplitter parity |
| T012 | Sentenizer trivial rules | T011 | Done | empty/no-space/lower/delimiter parity |
| T013 | Sentenizer sokr and initials rules | T012 | Planned | сокращения/инициалы parity |
| T014 | Sentenizer quote/bracket bound rules | T012 | Planned | кавычки/скобки на границах parity |
| T015 | Sentenizer bullet and dash rules | T012 | Planned | списки/тире-диалог parity |
| T016 | Sentenizer unit parity and post-processing API parity | T013,T014,T015 | Planned | unit parity + `strip()` + API-контракты соблюдены |
| T017 | Sentenizer integration corpus suite | T016,T002 | Planned | quick/full corpus режимы для sentenize |
| T018 | Differential/fuzz/benchmark hardening | T010,T017 | Planned | differential + fuzz + benchmark внедрены |

## Release gate checklist

- [ ] Все задачи `T001-T018` в статусе `Done`.
- [ ] `go test ./...` стабилен.
- [ ] `go vet ./...` стабилен.
- [ ] `golangci-lint run ./...` стабилен.
- [ ] Quick corpus parity в CI стабилен.
- [ ] Full corpus parity успешно завершен вручную/ночным прогоном.
- [ ] Все известные расхождения документированы.
