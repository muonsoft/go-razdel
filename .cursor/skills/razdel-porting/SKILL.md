---
name: razdel-porting
description: Port natasha/razdel behavior to Go with parity checks against upstream. Use when implementing tokenizer logic, validating edge cases, and documenting compatibility decisions.
---

# Razdel Porting

## Goal

Сохранять поведенческую совместимость Go-реализации с `third_party/razdel`.

## Workflow

1. Выбери небольшой срез поведения (один кейс или узкий набор кейсов).
2. Найди эталон в upstream (`third_party/razdel`) через тест или воспроизводимый пример.
3. Реализуй изменение в Go минимальным инкрементом.
4. Добавь/обнови тесты в Go, чтобы закрепить parity.
5. Зафиксируй отклонения и допущения в `README.md` или notes-файле.

## Parity Checklist

- [ ] Есть ссылка на upstream-кейс (тест или кодовый пример).
- [ ] Go-тест воспроизводит ожидаемое поведение.
- [ ] Краевые случаи Unicode и пунктуации проверены, если релевантны.
- [ ] Решение не изменяет `third_party/razdel`.

## Non-goals

- Не выполнять массовый перенос без тестового покрытия.
- Не обновлять submodule в рамках feature-изменений.
