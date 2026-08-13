# go-razdel

[English](README.md) | **Русский**

Сегментация русского текста на слова и предложения по правилам — порт [natasha/razdel](https://github.com/natasha/razdel) на Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/go-razdel.svg)](https://pkg.go.dev/github.com/muonsoft/go-razdel)
[![CI](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml/badge.svg)](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml)
[![Corpus](https://github.com/muonsoft/go-razdel/actions/workflows/corpus.yml/badge.svg)](https://github.com/muonsoft/go-razdel/actions/workflows/corpus.yml)
![Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/go-razdel?logo=go)
![Release](https://img.shields.io/github/v/release/muonsoft/go-razdel)
![License: MIT](https://img.shields.io/github/license/muonsoft/go-razdel)
[![Upstream](https://img.shields.io/badge/upstream-natasha%2Frazdel-3776AB?logo=python&logoColor=white)](https://github.com/natasha/razdel)
![Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)

`go-razdel` делит русский текст на **токены** и **предложения** тем же набором правил, что и Python `razdel`: сокращения (`т.д.`, инициалы), дроби и десятичные (`0.5`, `50/64`), кавычки, списки и тире.

Без моделей, без CGO, без словарей для скачивания. Нужен **Go 1.24+**.

Правила заточены под **новости и художественную прозу** (те же домены, что у оригинала). На соцсетях, научных статьях и юридических текстах качество может быть хуже — это ожидаемо.

## Установка

```bash
go get github.com/muonsoft/go-razdel
```

## Использование

Те же примеры, что в [natasha/razdel](https://github.com/natasha/razdel): рядом с Python-версией их удобно сравнивать глазами.

### Токены

```go
package main

import (
	"fmt"
	"strings"

	"github.com/muonsoft/go-razdel"
)

func main() {
	text := "Кружка-термос на 0.5л (50/64 см³, 516;...)"
	var parts []string
	for _, tok := range razdel.Tokenize(text) {
		parts = append(parts, tok.Text)
	}
	fmt.Println(strings.Join(parts, " | "))
}
```

```
Кружка-термос | на | 0.5 | л | ( | 50/64 | см³ | , | 516 | ; | ... | )
```

### Предложения

```go
text := `
- "Так в чем же дело?" - "Не ра-ду-ют".
И т. д. и т. п. В общем, вся газета
`

for _, sent := range razdel.Sentenize(text) {
	fmt.Println(sent.Text)
}
```

```
- "Так в чем же дело?"
- "Не ра-ду-ют".
И т. д. и т. п.
В общем, вся газета
```

Запускаемые копии этих сниппетов — в [`example_test.go`](example_test.go) и на [pkg.go.dev](https://pkg.go.dev/github.com/muonsoft/go-razdel).

## API

```go
tokens := razdel.Tokenize(text)   // []Token
sents  := razdel.Sentenize(text)  // []Sentence
```

У каждого элемента есть `Text` и полуинтервал `[Start, End)` в исходной строке (`Span` встроен, поэтому работает `tok.Start`):

```go
type Span struct {
	Start int // байтовый оффсет UTF-8, включительно
	End   int // байтовый оффсет UTF-8, исключительно
}

type Token struct {
	Span
	Text string // всегда равно text[Start:End]
}

type Sentence struct {
	Span
	Text string
}
```

Пустой ввод (и строка из одних пробелов для `Sentenize`) даёт nil-срез. Функции не возвращают `error` и не паникуют на обычном тексте, в том числе на невалидном UTF-8.

`Start`/`End` — это **байты UTF-8**, а не Unicode-кодпоинты. В тех же единицах работают `len` и срез `s[i:j]`, поэтому исходную строку можно резать напрямую:

```go
text := "a ж" // 4 байта, 3 руны: 'a' — 1 байт, 'ж' — 2
toks := razdel.Tokenize(text)
// "a" [0:1]
// "ж" [2:4]
fmt.Println(text[toks[1].Start:toks[1].End] == "ж") // true
```

Python `razdel` считает кодпоинты, поэтому числовые смещения на кириллице часто не совпадают, даже если **тексты** токенов одинаковые. Сравнивать реализации нужно по `Text`, а не по «сырым» индексам.

Смещения предложений указывают на **обрезанный** фрагмент (крайние пробелы сняты, как Python `chunk.strip()`).

## Совместимость с Python razdel

**Тексты** токенов и предложений следуют закреплённому [natasha/razdel](https://github.com/natasha/razdel) (`third_party/razdel`, коммит `668dbe191a5cfd94bebf9155e2ffa5f94ff3fe33`). В CI это проверяется unit-кейсами upstream, быстрым корпусом и живым сравнением с Python.

Два сознательных отличия tokenize (см. также upstream [razdel#17](https://github.com/natasha/razdel/issues/17) и [razdel#2](https://github.com/natasha/razdel/issues/2)):

| Вход | Python `razdel` | `go-razdel` |
| --- | --- | --- |
| `:-)` `;-)` `=-)` | `:`, `-`, `)` | один токен `:-)` |
| `✅Сдается` | один токен | `✅`, `Сдается` |
| `счетчики💰` | один токен | `счетчики`, `💰` |

Публичный API (`Tokenize`, `Sentenize`, `Token`, `Sentence`, `Span`, байтовые смещения) в линейке 0.x заморожен по смыслу: ломающие изменения идут с bump **minor** (0.2.0) и записью в [`CHANGELOG.md`](CHANGELOG.md). Go modules не обещают совместимость до v1.0.0.

## Документация

| Документ | О чём |
| --- | --- |
| [`README.md`](README.md) | Та же инструкция на английском |
| [pkg.go.dev/github.com/muonsoft/go-razdel](https://pkg.go.dev/github.com/muonsoft/go-razdel) | Сгенерированный API и примеры |
| [`docs/contracts.md`](docs/contracts.md) | Смещения, пустой ввод, невалидный UTF-8, известные отклонения |
| [`CHANGELOG.md`](CHANGELOG.md) | Пользовательские изменения |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Тесты, parity, релизы |
| [`docs/README.md`](docs/README.md) | Полный указатель документации |

## Участие

Issue и pull request приветствуются. Тесты, правила parity и кнопка релиза — в [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Лицензия

[MIT](LICENSE). Правила сегментации и списки сокращений взяты из [natasha/razdel](https://github.com/natasha/razdel) (MIT, Copyright 2017). Атрибуция: [NOTICE](NOTICE); дерево оригинала: `third_party/razdel`.
