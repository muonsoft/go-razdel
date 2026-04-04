# Upstream test inventory (`third_party/razdel`)

## Источники, просмотренные для планирования

- `razdel/tests/test_tokenize.py`
- `razdel/tests/test_sentenize.py`
- `razdel/tests/common.py`
- `razdel/tests/conftest.py`
- `razdel/tests/partition.py`
- `razdel/tests/gen.py`
- `razdel/tests/ctl.py`
- `razdel/tests/data/tokens.txt`
- `razdel/tests/data/sents.txt`

## Что тестирует upstream

### 1) `test_tokenize.py`

- `test_unit` — параметризованные unit-кейсы через `parse_partitions(...)`.
- `test_int` — интеграционные кейсы из `data/tokens.txt`.
- Механизм `--int N`: случайная выборка `N` строк с фиксированным `seed(1)`.

Покрываемые классы кейсов:
- дефисы/подчеркивания (`что-то`, `К_тому_же`);
- серия пунктуации (`...`, `:)))`, `: )||,`-подобные кейсы);
- числа (`1,5`, `1/2`);
- mixed scripts / unicode (`mβж`, `Δσ`);
- краевые пустые строки.

### 2) `test_sentenize.py`

- `test_unit` — вручную подобранные boundary-кейсы.
- `test_int` — интеграционные кейсы из `data/sents.txt`.
- Аналогичный режим `--int N`.

Покрываемые классы кейсов:
- базовые разделители (`.?!…`);
- аббревиатуры и сокращения (`т. д.`, `к.п.н.`, `т.е.`, `т.н.`);
- кавычки/скобки на границе предложения;
- тире-диалоги;
- списки/буллеты (`4.`, `IV.`, `§2.`, `8.1.`, `2)`).

## Масштаб интеграционных корпусов

Подсчет строк:

- `tokens.txt`: **208,995** строк.
- `sents.txt`: **48,735** строк.

Вывод: полный прогон на каждом PR может быть тяжелым; нужен режимы:
- быстрый deterministic sample для PR/CI;
- полный прогон по команде/ночной pipeline.

## Наблюдения по архитектуре upstream (для parity)

- Токенизация основана на:
  - атомизации по regex (`ATOM`),
  - `TokenSplitter`,
  - последовательных rule-объектах (dash/underscore/float/fraction/punct/other/yahoo).
- Сегментация предложений основана на:
  - `SentSplitter` + `SentSplit` вычисляемые свойства,
  - наборе правил (trivial/sokr/bound/list/dash),
  - post-step `strip()`.

## План тестовой миграции из upstream в Go

1. Повторить формат partition-тестов (`|` как граница сегмента).
2. Перенести unit-кейсы из `test_tokenize.py` и `test_sentenize.py` как go-table tests.
3. Поддержать интеграционный corpus-runner:
   - deterministic sample;
   - optional full corpus.
4. Добавить differential-тесты (Go vs upstream Python) для sampled кейсов.
