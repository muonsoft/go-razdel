# Changelog

Все заметные изменения `go-razdel` описываются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.1.0/), версии следуют [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Публичный API в линейке 0.x (`Tokenize`, `Sentenize`, `Span`, `Token`, `Sentence` и байтовые смещения Variant A) заморожен по смыслу: ломающие изменения сопровождаются bump minor и записью здесь. Go modules до v1.0.0 compatibility guarantee не дают.

## [Unreleased]

### Добавлено

- Первый публичный релиз библиотеки: `Tokenize` и `Sentenize` с поведенческой совместимостью с pinned [natasha/razdel](https://github.com/natasha/razdel) (`668dbe191a5cfd94bebf9155e2ffa5f94ff3fe33`).
- Контракт смещений Variant A: `Start`/`End` в байтах UTF-8, полуинтервал `[Start, End)`, инвариант `Text == text[Start:End]`.
- Unit-parity по upstream `UNIT`, quick corpus в CI, полный `tokens.txt` / `sents.txt` по флагу и в ночном workflow.
- Differential-тесты Go vs Python, fuzz по инвариантам смещений, бенчмарки.
- GitHub Actions workflow **Release**: кнопка в Actions публикует GitHub Release и переносит `[Unreleased]` из этого файла в версионную секцию.

### Изменено

- Контролируемые отклонения от Python tokenize (IGO-47): склейка смайлов `:-)` / `;-)` / `=-)` и отрыв emoji от соседних букв. См. README.
- Для невалидного UTF-8 зафиксировано фактическое поведение `Tokenize`/`Sentenize` (нет panic; Python parity для `Sentenize` не обещается).

### Известные ограничения

- Числовые `Start`/`End` в Go и Python не обязаны совпадать на не-ASCII тексте (байты vs кодпоинты).
- Версия 0.x: совместимость для потребителей Go-модуля не гарантируется языком до v1.0.0.
