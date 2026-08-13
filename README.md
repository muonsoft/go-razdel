# go-razdel

**English** | [Русский](README.ru.md)

Rule-based Russian word and sentence tokenizer — a Go port of [natasha/razdel](https://github.com/natasha/razdel).

[![Go Reference](https://pkg.go.dev/badge/github.com/muonsoft/go-razdel.svg)](https://pkg.go.dev/github.com/muonsoft/go-razdel)
[![CI](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml/badge.svg)](https://github.com/muonsoft/go-razdel/actions/workflows/ci.yml)
[![Corpus](https://github.com/muonsoft/go-razdel/actions/workflows/corpus.yml/badge.svg)](https://github.com/muonsoft/go-razdel/actions/workflows/corpus.yml)
![Go version](https://img.shields.io/github/go-mod/go-version/muonsoft/go-razdel?logo=go)
![Release](https://img.shields.io/github/v/release/muonsoft/go-razdel)
![License: MIT](https://img.shields.io/github/license/muonsoft/go-razdel)
[![Upstream](https://img.shields.io/badge/upstream-natasha%2Frazdel-3776AB?logo=python&logoColor=white)](https://github.com/natasha/razdel)
![Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)

`go-razdel` splits Russian text into **tokens** and **sentences** with the same rule set as Python `razdel`: abbreviations (`т.д.`, initials), decimals (`0.5`, `50/64`), quotes, lists, and dashes.

No models, no CGO, no extra dictionaries to download. Requires **Go 1.24+**.

Rules are tuned for **news and fiction** (the same domains as upstream). On social media, scientific papers, or legal text the result may be worse — that is expected.

## Install

```bash
go get github.com/muonsoft/go-razdel
```

## Usage

Same examples as [natasha/razdel](https://github.com/natasha/razdel), so a Python user can compare the two by eye.

### Tokens

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

### Sentences

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

Runnable copies of these snippets live in [`example_test.go`](example_test.go) and on [pkg.go.dev](https://pkg.go.dev/github.com/muonsoft/go-razdel).

## API

```go
tokens := razdel.Tokenize(text)   // []Token
sents  := razdel.Sentenize(text)  // []Sentence
```

Each item has `Text` plus a half-open span `[Start, End)` into the original string (`Span` is embedded, so `tok.Start` works):

```go
type Span struct {
	Start int // UTF-8 byte offset, inclusive
	End   int // UTF-8 byte offset, exclusive
}

type Token struct {
	Span
	Text string // always equal to text[Start:End]
}

type Sentence struct {
	Span
	Text string
}
```

Empty input (and whitespace-only input for `Sentenize`) returns a nil slice. Neither function returns an error or panics on ordinary text, including invalid UTF-8.

`Start`/`End` are **UTF-8 bytes**, not Unicode code points. That is the unit Go uses for `len` and `s[i:j]`, so you can slice the original string directly:

```go
text := "a ж" // 4 bytes, 3 runes: 'a' is 1 byte, 'ж' is 2
toks := razdel.Tokenize(text)
// "a" [0:1]
// "ж" [2:4]
fmt.Println(text[toks[1].Start:toks[1].End] == "ж") // true
```

Python `razdel` counts code points, so numeric offsets often differ on Cyrillic even when the token **texts** match. Compare implementations by `Text`, not by raw indexes.

Sentence spans point at the **trimmed** slice (leading/trailing whitespace is dropped, same as Python `chunk.strip()`).

## Compatibility with Python razdel

Token and sentence **texts** follow pinned [natasha/razdel](https://github.com/natasha/razdel) (`third_party/razdel`, commit `668dbe191a5cfd94bebf9155e2ffa5f94ff3fe33`), checked in CI against upstream unit cases, a quick corpus, and a live Python differential.

Two intentional tokenize differences (also discussed upstream as [razdel#17](https://github.com/natasha/razdel/issues/17) and [razdel#2](https://github.com/natasha/razdel/issues/2)):

| Input | Python `razdel` | `go-razdel` |
| --- | --- | --- |
| `:-)` `;-)` `=-)` | `:`, `-`, `)` | one token `:-)` |
| `✅Сдается` | one token | `✅`, `Сдается` |
| `счетчики💰` | one token | `счетчики`, `💰` |

The public API (`Tokenize`, `Sentenize`, `Token`, `Sentence`, `Span`, byte offsets) is frozen in meaning for 0.x: breaking changes bump **minor** (0.2.0) and are listed in [`CHANGELOG.md`](CHANGELOG.md). Go modules do not promise compatibility until v1.0.0.

## Documentation

| Document | What it covers |
| --- | --- |
| [`README.ru.md`](README.ru.md) | Same guide in Russian |
| [pkg.go.dev/github.com/muonsoft/go-razdel](https://pkg.go.dev/github.com/muonsoft/go-razdel) | Generated API reference and examples |
| [`docs/contracts.md`](docs/contracts.md) | Offsets, empty input, invalid UTF-8, known deviations |
| [`CHANGELOG.md`](CHANGELOG.md) | User-visible changes |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Tests, parity checks, releases |
| [`docs/README.md`](docs/README.md) | Full documentation index |

## Contributing

Issues and pull requests are welcome. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for tests, parity rules, and the release button.

## License

[MIT](LICENSE). Segmentation rules and abbreviation lists are derived from [natasha/razdel](https://github.com/natasha/razdel) (MIT, Copyright 2017). Attribution: [NOTICE](NOTICE); upstream tree: `third_party/razdel`.
