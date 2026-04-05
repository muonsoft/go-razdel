# go-razdel

Go-порт библиотеки [natasha/razdel](https://github.com/natasha/razdel).

## Status

Проект в ранней стадии:
- инициализирован Go module;
- upstream подключен как git submodule;
- публичный API и типы соответствуют `docs/contracts.md`;
- `Tokenize` реализован с ориентацией на upstream; `Sentenize` применяет trivial-слой join-правил upstream (остальные правила — в следующих задачах).

Реализация алгоритмов и parity с upstream выполняются поэтапно.

## Repository Layout

- `go.mod` — модуль `github.com/muonsoft/go-razdel`.
- `third_party/razdel` — upstream submodule с эталонной реализацией.
- `AGENTS.md` — правила для агентной работы в репозитории.
- `docs/contracts.md` — контрактные требования v0 (включая байтовые смещения UTF-8).
- `.github/workflows/ci.yml` — CI: `go test`, `go vet`, `golangci-lint`.
- `.golangci.yml` — конфигурация линтера (schema v2).

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

Локально (при установленном [golangci-lint](https://golangci-lint.run/)):

```bash
golangci-lint run ./...
```

## Upstream Pinning

Submodule всегда фиксируется на конкретный commit upstream. Обновлять `third_party/razdel` нужно только отдельной, явной задачей и отдельным коммитом, чтобы изменения поведения были полностью трассируемы.

## Contracts

Верхнеуровневые контрактные требования описаны в `docs/contracts.md`.
