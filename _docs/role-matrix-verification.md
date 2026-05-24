# Ролевая матрица — Чеклист верификации

> Дата последней проверки: —
> Проверял: —
> Окружение: local dev / staging

## Легенда

| Символ | Значение                       |
| ------ | ------------------------------ |
| ✅     | Проверено, работает            |
| ❌     | Проверено, ОШИБКА — нужен фикс |
| ⏳     | Не проверено                   |
| N/A    | Не применимо без полного стека |

---

## Часть 1 — Автоматизированная (Playwright, `pnpm e2e`)

Запуск: `NEXTAUTH_SECRET=change-me-in-production pnpm e2e`

| Тест                  | Описание                                                        | Результат |
| --------------------- | --------------------------------------------------------------- | --------- |
| `redirects.spec.ts`   | public → /login при попытке войти в защищённые разделы          | ⏳        |
| `cabinet.spec.ts`     | borrower: все страницы ЛК работают                              | ⏳        |
| `employee.spec.ts`    | employee: экспертиза и аналитика работают                       | ⏳        |
| `role-matrix.spec.ts` | cross-role: borrower → /403 при попытке /expertise и /analytics | ⏳        |
| `role-matrix.spec.ts` | expert: доступ к /expertise, /analytics                         | ⏳        |
| `role-matrix.spec.ts` | admin: доступ ко всем разделам                                  | ⏳        |
| `role-matrix.spec.ts` | /403 страница рендерится корректно                              | ⏳        |

---

## Часть 2 — Frontend (Next.js middleware), ручная проверка

Запуск dev-сервера: `E2E_TEST=true NEXTAUTH_SECRET=... pnpm dev`

### public (без сессии)

| URL           | Ожидаемый результат                        | Результат |
| ------------- | ------------------------------------------ | --------- |
| `/`           | 200, главная страница                      | ⏳        |
| `/programs`   | 200, список программ                       | ⏳        |
| `/calculator` | 200, калькулятор                           | ⏳        |
| `/map`        | 200, публичная карта                       | ⏳        |
| `/faq`        | 200, FAQ                                   | ⏳        |
| `/contacts`   | 200, контакты                              | ⏳        |
| `/cabinet`    | redirect → `/login?callbackUrl=/cabinet`   | ⏳        |
| `/expertise`  | redirect → `/login?callbackUrl=/expertise` | ⏳        |
| `/analytics`  | redirect → `/login?callbackUrl=/analytics` | ⏳        |

### borrower

| URL                                       | Ожидаемый результат  | Результат |
| ----------------------------------------- | -------------------- | --------- |
| `/cabinet`                                | 200, дашборд ЛК      | ⏳        |
| `/cabinet/schedule`                       | 200, график платежей | ⏳        |
| `/cabinet/documents`                      | 200, документы       | ⏳        |
| `/cabinet/tickets`                        | 200, обращения       | ⏳        |
| `/cabinet/notifications`                  | 200, уведомления     | ⏳        |
| `/expertise`                              | redirect → `/403`    | ⏳        |
| `/analytics`                              | redirect → `/403`    | ⏳        |
| Кнопки «Аналитика», «Экспертиза» в хедере | Не видны             | ⏳        |

### employee

| URL                                       | Ожидаемый результат              | Результат |
| ----------------------------------------- | -------------------------------- | --------- |
| `/cabinet`                                | 200 (employee может смотреть ЛК) | ⏳        |
| `/expertise`                              | 200, очередь заявок              | ⏳        |
| `/analytics`                              | 200, Grafana-дашборд             | ⏳        |
| Кнопки «Аналитика», «Экспертиза» в хедере | Видны                            | ⏳        |

### expert

| URL          | Ожидаемый результат | Результат |
| ------------ | ------------------- | --------- |
| `/cabinet`   | 200                 | ⏳        |
| `/expertise` | 200                 | ⏳        |
| `/analytics` | 200                 | ⏳        |

### admin

| URL          | Ожидаемый результат | Результат |
| ------------ | ------------------- | --------- |
| `/cabinet`   | 200                 | ⏳        |
| `/expertise` | 200                 | ⏳        |
| `/analytics` | 200                 | ⏳        |

---

## Часть 3 — Go API (core-api :8080), ручная проверка

> Требует запущенный стек: `docker compose up -d`
> Используй seed-пользователей из `apps/api/core-api/migrations/seed/`.

### Базовые команды для проверки

```bash
# Получить токен для роли (через Keycloak)
TOKEN=$(curl -s -X POST http://localhost:8180/realms/pkt/protocol/openid-connect/token \
  -d "client_id=web&grant_type=password&username=borrower@pkt.kz&password=test123" \
  | jq -r .access_token)

# Проверить эндпоинт
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/programs
```

### public (без токена)

| Эндпоинт                    | Метод | Ожидание | Результат |
| --------------------------- | ----- | -------- | --------- |
| `GET /api/v1/health`        | GET   | 200      | ⏳        |
| `GET /api/v1/programs`      | GET   | 200      | ⏳        |
| `GET /api/v1/programs/{id}` | GET   | 200      | ⏳        |
| `POST /api/v1/applications` | POST  | 401      | ⏳        |
| `GET /api/v1/me`            | GET   | 401      | ⏳        |
| `GET /api/v1/notifications` | GET   | 401      | ⏳        |
| `GET /api/v1/admin/users`   | GET   | 401      | ⏳        |

### borrower

| Эндпоинт                    | Метод | Ожидание          | Результат |
| --------------------------- | ----- | ----------------- | --------- |
| `GET /api/v1/me`            | GET   | 200               | ⏳        |
| `GET /api/v1/me/borrower`   | GET   | 200               | ⏳        |
| `POST /api/v1/applications` | POST  | 201               | ⏳        |
| `GET /api/v1/applications`  | GET   | 200 (только свои) | ⏳        |
| `GET /api/v1/tickets`       | GET   | 200 (только свои) | ⏳        |
| `POST /api/v1/tickets`      | POST  | 201               | ⏳        |
| `GET /api/v1/notifications` | GET   | 200               | ⏳        |
| `GET /api/v1/users`         | GET   | **403**           | ⏳        |
| `POST /api/v1/programs`     | POST  | **403**           | ⏳        |
| `GET /api/v1/admin/users`   | GET   | **403**           | ⏳        |

### employee

| Эндпоинт                                 | Метод | Ожидание         | Результат |
| ---------------------------------------- | ----- | ---------------- | --------- |
| `GET /api/v1/users`                      | GET   | 200              | ⏳        |
| `GET /api/v1/applications`               | GET   | 200 (все заявки) | ⏳        |
| `PATCH /api/v1/applications/{id}/status` | PATCH | 200              | ⏳        |
| `GET /api/v1/tickets`                    | GET   | 200 (все тикеты) | ⏳        |
| `POST /api/v1/applications`              | POST  | **403**          | ⏳        |
| `POST /api/v1/programs`                  | POST  | **403**          | ⏳        |
| `GET /api/v1/admin/users`                | GET   | **403**          | ⏳        |

### admin

| Эндпоинт                       | Метод  | Ожидание | Результат |
| ------------------------------ | ------ | -------- | --------- |
| `GET /api/v1/admin/users`      | GET    | 200      | ⏳        |
| `POST /api/v1/programs`        | POST   | 201      | ⏳        |
| `DELETE /api/v1/programs/{id}` | DELETE | 200      | ⏳        |
| `POST /api/v1/applications`    | POST   | 201      | ⏳        |

---

## Часть 4 — expertise-svc (:8081)

### employee vs expert vs borrower

| Эндпоинт                                              | Роль     | Ожидание        | Результат |
| ----------------------------------------------------- | -------- | --------------- | --------- |
| `GET /api/v1/expertise/queue`                         | employee | 200             | ⏳        |
| `GET /api/v1/expertise/queue`                         | expert   | 200             | ⏳        |
| `GET /api/v1/expertise/queue`                         | borrower | **403**         | ⏳        |
| `POST /api/v1/expertise/applications/{id}/transition` | employee | 200 (свой этап) | ⏳        |
| `POST /api/v1/expertise/applications/{id}/vote`       | employee | **403**         | ⏳        |
| `POST /api/v1/expertise/applications/{id}/vote`       | expert   | 200             | ⏳        |
| `POST /api/v1/expertise/committee/applications`       | employee | **403**         | ⏳        |
| `POST /api/v1/expertise/committee/applications`       | expert   | 200             | ⏳        |

---

## Часть 5 — UI/UX элементы по ролям

Проверяется визуально в браузере.

| Элемент                             | Видят            | Не видят         | Результат |
| ----------------------------------- | ---------------- | ---------------- | --------- |
| Кнопка «Подать заявку» в хедере     | public, borrower | —                | ⏳        |
| Ссылка «Личный кабинет» в хедере    | borrower+        | public           | ⏳        |
| Ссылка «Аналитика» в хедере         | employee+        | public, borrower | ⏳        |
| Ссылка «Экспертиза» в хедере        | employee+        | public, borrower | ⏳        |
| Инлайн-редактор контента (programs) | employee+        | остальные        | ⏳        |
| Финансовая карточка в popup карты   | employee+        | public, borrower | ⏳        |
| Слои тепловой карты портфеля        | employee+        | public, borrower | ⏳        |
| Кнопка «Добавить FAQ»               | employee+        | остальные        | ⏳        |
| Чат-виджет                          | все роли         | —                | ⏳        |

---

## Итог

| Категория           | Всего | ✅  | ❌  | Замечания |
| ------------------- | ----- | --- | --- | --------- |
| Playwright E2E      | —     | —   | —   |           |
| Frontend middleware | —     | —   | —   |           |
| core-api RBAC       | —     | —   | —   |           |
| expertise-svc RBAC  | —     | —   | —   |           |
| UI/UX элементы      | —     | —   | —   |           |

**Статус задачи 5.5:** ⏳ В процессе
