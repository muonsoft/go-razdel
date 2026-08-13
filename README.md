# go-razdel

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/go-razdel.svg)](https://pkg.go.dev/github.com/muonsoft/go-razdel)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/go-razdel)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/muonsoft/go-razdel)
![GitHub](https://img.shields.io/github/license/muonsoft/go-razdel)
[![tests](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml/badge.svg)](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml)

`go-razdel` — Go-порт библиотеки сегментации текста **razdel** с фокусом на поведенческую совместимость с Python-оригиналом.

Первый фиксированный релиз: **v0.1.0**. История изменений — [`CHANGELOG.md`](CHANGELOG.md).

## Основной проект

Эталонная реализация, на которую ориентируется этот репозиторий:

- [natasha/razdel (GitHub)](https://github.com/natasha/razdel)
- [razdel (PyPI)](https://pypi.org/project/razdel/)

Upstream закреплён git-submodule `third_party/razdel` (MIT, Copyright 2017). Правила сегментации и списки сокращений в Go-коде производны от этого дерева; см. также `LICENSE`.

## Возможности

- Публичный API:
  - `Tokenize(text string) []Token`
  - `Sentenize(text string) []Sentence`
- Контракт смещений: `Start`/`End` в **байтах UTF-8**, полуинтервал `[Start, End)`.
- `Tokenize` и `Sentenize` возвращают `Text == text[Start:End]` для каждого элемента.
- Проверки parity с upstream:
  - unit-набор для `Tokenize` и `Sentenize`;
  - differential-тесты (Go vs Python) для быстрых выборок токенов и предложений.
- Фаззинг и инварианты для оффсетов.

## Совместимость (v0.1)

Публичный API v0.1 (`Tokenize`, `Sentenize`, `Span`, `Token`, `Sentence`) и байтовые смещения Variant A **заморожены по смыслу**. Ломающие изменения в 0.x идут с bump minor (например 0.2.0) и записью в CHANGELOG.

Это линейка **0.x**: [Go modules](https://go.dev/blog/v2-go-modules) не дают compatibility guarantee до v1.0.0. После soak планируется v1.0.0 с той же семантикой, если API не придётся ломать.

Минимальная версия языка: **Go 1.24**.

## Ограничения и важные отличия

- Числовые смещения в Go и Python могут отличаться на Unicode-тексте, потому что:
  - в Go API фиксирует **байтовые** индексы;
  - в Python upstream работает с индексами `str` (кодпоинты).
- Parity корректно сравнивать по последовательности `Text`, а не по «сырым» индексам без пересчета.

### Невалидный UTF-8

Python `str` не содержит невалидный UTF-8; Go-строка может.

- `Tokenize` не паникует. Каждый невалидный байт — атом `OTHER`. Join-правила те же, что для прочих `OTHER`: соседние `OTHER`/`RU`/`LAT` склеиваются, поэтому `a\xffb` даёт один токен `"a\xffb"`, а не три.
- `Sentenize` не паникует, инвариант смещений сохраняется. Parity с Python для такого ввода не обещается.
- Внутренний `ByteSpans` паникует только если stripped-чанк не найден в исходной строке (программная ошибка пайплайна, не пользовательский ввод).

### Токенизация: смайлы и emoji (IGO-47)

Зафиксировано **контролируемое отклонение** от pinned `razdel.segmenters.tokenize` (см. [upstream #17](https://github.com/natasha/razdel/issues/17), [#2](https://github.com/natasha/razdel/issues/2)):

- Смайлы вида `:-)`, `;-)`, `=-)` склеиваются в один токен (lookahead к следующему пунктуационному атому), чтобы не резать на `:`, `-`, `)`.
- Эмодзи (по эвристике, близкой к Unicode Extended Pictographic) **отделяются** от соседних кириллических/латинских букв: например `✅Сдается` → `✅`, `Сдается`; `счетчики💰` → `счетчики`, `💰`.

Инвариант смещений сохраняется: для каждого токена `Text == text[Start:End]` (байтовые индексы UTF-8).

Детали контрактов и инвариантов: `docs/contracts.md`.
Базовые метрики и процесс контроля регрессий: `docs/performance.md`.

## Установка

```bash
go get github.com/muonsoft/go-razdel
```

Для разработки в этом репозитории (нужен upstream submodule):

```bash
git clone git@github.com:muonsoft/go-razdel.git
cd go-razdel
git submodule update --init --recursive
```

## Быстрый пример

```go
package main

import (
	"fmt"

	"github.com/muonsoft/go-razdel"
)

func main() {
	text := "Привет, мир! Это тест."

	for _, tok := range razdel.Tokenize(text) {
		fmt.Printf("TOKEN  [%d:%d] %q\n", tok.Start, tok.End, tok.Text)
	}

	for _, sent := range razdel.Sentenize(text) {
		fmt.Printf("SENT   [%d:%d] %q\n", sent.Start, sent.End, sent.Text)
	}
}
```

## Структуры данных

- `Span`:
  - `Start int`
  - `End int`
- `Token`:
  - `Span`
  - `Text string`
- `Sentence`:
  - `Span`
  - `Text string`

## Проверка и разработка

Минимальный локальный набор:

```bash
go test ./...
go vet ./...
golangci-lint run ./...
```

Differential-тесты с Python требуют:

- `python3` в `PATH`;
- доступный submodule `third_party/razdel`;
- импорт `razdel` через `PYTHONPATH` (настраивается тестами автоматически).

В CI differential **обязан** выполниться (`RAZDEL_DIFFERENTIAL_PYTHON=1`: skip превращается в fail). Локально skip сохраняется, если Python недоступен.

Отключить differential-тесты:

```bash
RAZDEL_DIFFERENTIAL_PYTHON=0 go test ./...
```

Полный корпус (~209k tokenize / ~49k sentenize) — вручную или ночной workflow `.github/workflows/corpus.yml`:

```bash
RAZDEL_TOKENIZE_INTEGRATION_FULL=1 go test ./internal/tokenize/ -run TestIntegration_tokenize_full_corpus -count=1
RAZDEL_SENTENIZE_INTEGRATION_FULL=1 go test ./internal/sentenize/ -run TestIntegration_sentenize_full_corpus -count=1
```

## Релизы

Релиз создаётся кнопкой в GitHub Actions (**Actions → Release → Run workflow**) с версией вида `v0.1.0`. Workflow также срабатывает на push тега `v*.*.*`.

Перед релизом добавляйте пользовательские изменения в секцию **`## [Unreleased]`** в [`CHANGELOG.md`](CHANGELOG.md). При публикации workflow:

1. прогоняет тесты (включая обязательный Python differential), vet и lint;
2. проверяет, что в `[Unreleased]` есть записи;
3. переносит `[Unreleased]` в секцию `[x.y.z]` с датой UTC;
4. оставляет пустой `[Unreleased]` для следующего цикла;
5. коммитит changelog в `main`;
6. публикует GitHub Release; текст релиза берётся из новой версионной секции.

Если секция для версии уже есть в `CHANGELOG.md`, workflow не дублирует её и только публикует релиз.

Запускайте **Run workflow** с ветки `main`.
