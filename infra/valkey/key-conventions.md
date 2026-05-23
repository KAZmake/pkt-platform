# Valkey — конвенции ключей PKT Platform

## Логические базы данных

| DB  | Назначение        | Политика     |
| --- | ----------------- | ------------ |
| 0   | Сессии            | noeviction   |
| 1   | Кэш API           | volatile-lru |
| 2   | Rate limiting     | volatile-lru |
| 3   | Distributed locks | volatile-lru |

---

## Префиксы ключей

### DB 0 — Сессии (`session:`)

| Ключ                        | Значение                                | TTL |
| --------------------------- | --------------------------------------- | --- |
| `session:{user_id}`         | JSON: `{keycloak_id, role, email, ...}` | 24h |
| `session:refresh:{user_id}` | Refresh token hash                      | 7d  |

### DB 1 — Кэш API (`cache:`)

| Ключ                        | Значение                           | TTL |
| --------------------------- | ---------------------------------- | --- |
| `cache:programs`            | JSON: список программ кредитования | 30m |
| `cache:program:{id}`        | JSON: одна программа               | 30m |
| `cache:loans:{borrower_id}` | JSON: займы заёмщика из 1С         | 20m |
| `cache:schedule:{loan_id}`  | JSON: график платежей              | 20m |
| `cache:debts:{borrower_id}` | JSON: задолженности                | 20m |
| `cache:collateral:{id}`     | JSON: данные залога                | 30m |
| `cache:news`                | JSON: список новостей (публичный)  | 15m |
| `cache:faq`                 | JSON: FAQ (публичный)              | 60m |

### DB 2 — Rate limiting (`ratelimit:`)

| Ключ                                  | Значение                                | TTL |
| ------------------------------------- | --------------------------------------- | --- |
| `ratelimit:ip:{ip}:{endpoint}`        | counter (INCR)                          | 60s |
| `ratelimit:user:{user_id}:{endpoint}` | counter (INCR)                          | 60s |
| `ratelimit:assistant:{ip}`            | counter — лимит 20 req/min для чат-бота | 60s |

### DB 3 — Distributed locks (`lock:`)

| Ключ                    | Значение                                | TTL |
| ----------------------- | --------------------------------------- | --- |
| `lock:sync:1c`          | `1` — синхронизация с 1С (один процесс) | 30s |
| `lock:application:{id}` | `{worker_id}` — обработка заявки        | 30s |

---

## Правила именования

- Разделитель: `:` (двоеточие)
- UUID без дефисов: `cache:loan:550e8400e29b41d4a716446655440000`
- Нижний регистр везде
- Версионирование кэша через суффикс: `cache:programs:v2` (при смене схемы)
