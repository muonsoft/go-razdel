# go-razdel

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/go-razdel.svg)](https://pkg.go.dev/github.com/muonsoft/go-razdel)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/go-razdel)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/muonsoft/go-razdel)
![GitHub](https://img.shields.io/github/license/muonsoft/go-razdel)
[![tests](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml/badge.svg)](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml)

`go-razdel` — Go-порт библиотеки сегментации текста **razdel** с фокусом на поведенческую совместимость с Python-оригиналом.

## Основной проект

Эталонная реализация, на которую ориентируется этот репозиторий:

- [natasha/razdel (GitHub)](https://github.com/natasha/razdel)
- [razdel (PyPI)](https://pypi.org/project/razdel/)

## Что уже реализовано

- Публичный API:
  - `Tokenize(text string) []Token`
  - `Sentenize(text string) []Sentence`
- Контракт смещений: `Start`/`End` в **байтах UTF-8**, полуинтервал `[Start, End)`.
- `Tokenize` и `Sentenize` возвращают `Text == text[Start:End]` для каждого элемента.
- Проверки parity с upstream:
  - unit-набор для `Sentenize` (кейсы из upstream фикстур);
  - differential-тесты (Go vs Python) для быстрых выборок токенов и предложений.
- Фаззинг и инварианты для оффсетов.

## Ограничения и важные отличия

- Числовые смещения в Go и Python могут отличаться на Unicode-тексте, потому что:
  - в Go API фиксирует **байтовые** индексы;
  - в Python upstream работает с индексами `str` (кодпоинты).
- Parity корректно сравнивать по последовательности `Text`, а не по «сырым» индексам без пересчета.
- Для невалидного UTF-8 `Tokenize` не паникует; для `Sentenize` parity с Python не гарантируется.

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

Отключить differential-тесты можно переменной:

```bash
RAZDEL_DIFFERENTIAL_PYTHON=0 go test ./...
```

## Статус проекта

Проект развивается инкрементально: поведение переносится небольшими шагами, а совместимость с upstream подтверждается тестами и фикстурами.
