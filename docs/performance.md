# Performance guide

Этот документ фиксирует базовые метрики производительности и процесс проверки регрессий для `go-razdel`.

## Scope

- Проверяем публичные API:
  - `Tokenize(text string) []Token`
  - `Sentenize(text string) []Sentence`
- Основные сценарии нагрузки:
  - `short_ASCII`
  - `long_Cyrillic`
  - `punct_heavy`
  - `unicode_mixed`

## Baseline (после оптимизаций hot path)

Команда:

```bash
go test -bench=. -benchmem ./...
```

Среда замера:

- OS: linux
- arch: amd64
- CPU: Intel Xeon (cloud VM)
- Go: 1.25.3

Результаты:

### Tokenize

| benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| short_ASCII | 873.2 | 1696 | 12 |
| long_Cyrillic | 8183 | 10816 | 43 |
| punct_heavy | 353487 | 720961 | 3046 |
| unicode_mixed | 40993 | 56768 | 282 |

### Sentenize

| benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| short_ASCII | 7085 | 2257 | 39 |
| long_Cyrillic | 19357 | 3915 | 45 |
| punct_heavy | 1887034 | 1771658 | 4816 |
| unicode_mixed | 92684 | 60648 | 247 |

## Regression check workflow

1. Снять свежие метрики:

```bash
go test -run=^$ -bench='Benchmark(Tokenize|Sentenize)' -benchmem -count=5 ./...
```

2. Сравнить с baseline из этого файла:
   - в первую очередь `ns/op` и `allocs/op`;
   - отдельно контролировать `punct_heavy`, так как это самый чувствительный сценарий.

3. Если рост метрик заметный, приложить профили:

```bash
go test -run=^$ -bench='BenchmarkTokenize/punct_heavy$' -benchmem -cpuprofile=tokenize.cpu -memprofile=tokenize.mem .
go test -run=^$ -bench='BenchmarkSentenize/punct_heavy$' -benchmem -cpuprofile=sentenize.cpu -memprofile=sentenize.mem .
go tool pprof -top ./go-razdel.test tokenize.cpu
go tool pprof -top -alloc_space ./go-razdel.test tokenize.mem
go tool pprof -top ./go-razdel.test sentenize.cpu
go tool pprof -top -alloc_space ./go-razdel.test sentenize.mem
```

## Policy for baseline updates

- Baseline обновляется только вместе с изменением производительности (осознанно).
- В PR нужно указывать:
  - какие функции ускорены/замедлены;
  - на каких сценариях;
  - какие компромиссы приняты (например, рост памяти ради скорости).
- Если поведение и API не менялись, но есть деградация >10% по `ns/op` или >15% по `allocs/op` на критичных бенчмарках, это нужно считать риском и явно разбирать в PR.
