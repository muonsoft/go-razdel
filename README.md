# go-razdel

Go-порт библиотеки [natasha/razdel](https://github.com/natasha/razdel).

## Status

Проект находится в стадии bootstrap/pre-implementation:
- инициализирован Go module;
- upstream подключен как git submodule;
- подготовлены базовые правила для дальнейшей итеративной разработки.

На этом этапе перенос алгоритмов из Python в Go не выполняется.

## Repository Layout

- `go.mod` — модуль `github.com/muonsoft/go-razdel`.
- `third_party/razdel` — upstream submodule с эталонной реализацией.
- `AGENTS.md` — правила для агентной работы в репозитории.

## Getting Started

```bash
git clone git@github.com:muonsoft/go-razdel.git
cd go-razdel
git submodule update --init --recursive
```

Проверка окружения Go:

```bash
go version
go list ./...
```

## Upstream Pinning

Submodule всегда фиксируется на конкретный commit upstream. Обновлять `third_party/razdel` нужно только отдельной, явной задачей и отдельным коммитом, чтобы изменения поведения были полностью трассируемы.
