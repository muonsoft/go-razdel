# T015 — Sentenizer bullet and dash rules

## Goal

Перенести правила для bullet/list элементов и тире-диалогов.

## Upstream references

- `third_party/razdel/razdel/segmenters/sentenize.py`:
  - `is_bullet`
  - `list_item`
  - `dash_right`
  - `ROMAN`, `BULLET_CHARS`, `BULLET_BOUNDS`, `BULLET_SIZE`
- `test_sentenize.py` cases:
  - `4.`, `IV.`, `§2.`, `8.1.`, `2)`, диалоговые тире

## Scope

### In scope
- Bullet эвристика (digits/roman/letter bounds).
- Ограничение длины `BULLET_SIZE`.
- Dash-right rule для продолжения с lowercase.

### Out of scope
- Общая нормализация списков за пределами upstream логики.

## Steps

1. Перенести `is_bullet` + `list_item`.
2. Перенести `dash_right`.
3. Добавить unit тесты на списки и диалоговые реплики.

## Tests

- Unit: numbered/roman/paragraph bullets.
- Unit: dialogue dash cases.
- Regression: не ломать quote/sokr решения.

## Acceptance criteria

- [ ] Bullet/dash поведение соответствует upstream.
- [ ] Кейсы из upstream UNIT покрыты.
- [ ] `go test ./...` проходит.
