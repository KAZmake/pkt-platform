# Ролевая матрица — PKT Platform

> Роль × Эндпоинт × Действие.
> Используется для настройки JWT middleware в Go-сервисах.

## Обозначения

| Символ | Значение |
|--------|---------|
| ✓ | Разрешено |
| ✗ | Запрещено (403) |
| — | Не применимо |
| ✓* | Только свои данные |

## Роли

| Роль | Код | Описание |
|------|-----|---------|
| Публичный посетитель | `public` | Без авторизации |
| Заёмщик | `borrower` | Авторизован, личный кабинет |
| Сотрудник | `employee` | Персонал, экспертиза (свой этап), аналитика |
| Эксперт / КК | `expert` | Экспертиза залогов, юр. проверка, голосование КК |
| Администратор | `admin` | Полный доступ |

---

## core-api (порт 8080)

### Health & Profile

| Эндпоинт | Метод | public | borrower | employee | expert | admin |
|----------|-------|--------|---------|---------|--------|-------|
| `/health` | GET | ✓ | ✓ | ✓ | ✓ | ✓ |
| `/profile` | GET | ✗ | ✓ | ✓ | ✓ | ✓ |
| `/profile` | PATCH | ✗ | ✓ | ✓ | ✓ | ✓ |

### Программы кредитования

| Эндпоинт | Метод | public | borrower | employee | expert | admin |
|----------|-------|--------|---------|---------|--------|-------|
| `/programs` | GET | ✓ | ✓ | ✓ | ✓ | ✓ |
| `/programs/{id}` | GET | ✓ | ✓ | ✓ | ✓ | ✓ |

> Создание/редактирование программ — через Directus CMS (employee+), не через API напрямую.

### Заявки

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/applications` | POST | ✗ | ✓ | ✗ | ✗ | ✓ | Только заёмщик подаёт |
| `/applications` | GET | ✗ | ✓* | ✓ | ✓ | ✓ | borrower — только свои |
| `/applications/{id}` | GET | ✗ | ✓* | ✓ | ✓ | ✓ | borrower — только свои |

### ЛК — Документы

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/lk/documents` | GET | ✗ | ✓* | ✓ | ✗ | ✓ | borrower — только свои |
| `/lk/documents/{id}/download` | GET | ✗ | ✓* | ✓ | ✗ | ✓ | presigned URL MinIO |

### ЛК — Кредиты и платежи

| Эндпоинт | Метод | public | borrower | employee | expert | admin |
|----------|-------|--------|---------|---------|--------|-------|
| `/lk/loans` | GET | ✗ | ✓* | ✓ | ✗ | ✓ |
| `/lk/loans/{id}/schedule` | GET | ✗ | ✓* | ✓ | ✗ | ✓ |

### ЛК — Обращения (тикеты)

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/lk/tickets` | GET | ✗ | ✓* | ✓ | ✗ | ✓ | borrower — только свои |
| `/lk/tickets` | POST | ✗ | ✓ | ✗ | ✗ | ✓ | Только заёмщик создаёт |
| `/lk/tickets/{id}` | GET | ✗ | ✓* | ✓ | ✗ | ✓ | |
| `/lk/tickets/{id}/messages` | POST | ✗ | ✓* | ✓ | ✗ | ✓ | employee — только assigned |

### Уведомления

| Эндпоинт | Метод | public | borrower | employee | expert | admin |
|----------|-------|--------|---------|---------|--------|-------|
| `/notifications` | GET | ✗ | ✓* | ✓* | ✓* | ✓ |
| `/notifications/{id}/read` | PATCH | ✗ | ✓* | ✓* | ✓* | ✓ |
| `/notifications/read-all` | PATCH | ✗ | ✓* | ✓* | ✓* | ✓ |

### Администрирование

| Эндпоинт | Метод | public | borrower | employee | expert | admin |
|----------|-------|--------|---------|---------|--------|-------|
| `/admin/users` | GET | ✗ | ✗ | ✗ | ✗ | ✓ |
| `/admin/users/{id}/role` | PATCH | ✗ | ✗ | ✗ | ✗ | ✓ |

---

## expertise-svc (порт 8081)

### Очередь заявок

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/expertise/queue` | GET | ✗ | ✗ | ✓* | ✓* | ✓ | Видят только свой этап |

### Карточка заявки

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/expertise/applications/{id}` | GET | ✗ | ✗ | ✓* | ✓* | ✓ | Свой этап |
| `/expertise/applications/{id}/general` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/financial` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/collaterals` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/collaterals` | POST | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/collaterals/{cid}/release` | POST | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/documents` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/documents` | POST | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/history` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | Только чтение |

### FSM — Переходы статусов

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/expertise/applications/{id}/transition` | POST | ✗ | ✗ | ✓* | ✓* | ✓ | Только свой этап FSM |

**Допустимые переходы по ролям:**

| Переход | Кто может |
|---------|----------|
| `received → primary_scoring` | employee |
| `primary_scoring → security_check` | employee |
| `security_check → collateral_expertise` | employee |
| `collateral_expertise → legal_check` | expert |
| `legal_check → credit_analysis` | expert |
| `credit_analysis → credit_committee` | employee |
| `credit_committee → approved/rejected` | expert (после голосования КК) |
| `* → revision` | employee/expert/admin |
| `revision → *` | employee (возврат на нужный этап) |
| `approved → documentation` | employee |
| `documentation → issued` | employee/admin |

### Залоги

| Эндпоинт | Метод | public | borrower | employee | expert | admin |
|----------|-------|--------|---------|---------|--------|-------|
| `/expertise/collaterals` | GET | ✗ | ✗ | ✓ | ✓ | ✓ |
| `/expertise/collaterals` | POST | ✗ | ✗ | ✓ | ✓ | ✓ |
| `/expertise/collaterals/{id}` | GET | ✗ | ✗ | ✓ | ✓ | ✓ |
| `/expertise/collaterals/{id}` | PATCH | ✗ | ✗ | ✓ | ✓ | ✓ |

### Заключения экспертов

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/expertise/applications/{id}/conclusions` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/expertise/applications/{id}/conclusions` | POST | ✗ | ✗ | ✓* | ✓* | ✓ | Только свой этап |

### Кредитный комитет

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/expertise/committee/applications` | GET | ✗ | ✗ | ✗ | ✓ | ✓ | Только статус credit_committee |
| `/expertise/applications/{id}/vote` | POST | ✗ | ✗ | ✗ | ✓ | ✓ | |
| `/expertise/applications/{id}/protocol` | POST | ✗ | ✗ | ✗ | ✓ | ✓ | Генерация PDF |

---

## sync-svc (порт 8082)

> Внутренний сервис. Клиенты (web/mobile) не обращаются напрямую — только через core-api/expertise-svc.

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/sync/trigger` | POST | ✗ | ✗ | ✗ | ✗ | ✓ | Ручной запуск синхронизации |
| `/sync/status` | GET | ✗ | ✗ | ✓ | ✗ | ✓ | |
| `/sync/loans` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/sync/loans/{id}` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/sync/loans/{id}/schedule` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/sync/loans/{id}/debts` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/sync/collateral-alerts` | GET | ✗ | ✗ | ✓ | ✓ | ✓ | |
| `/sync/upcoming-payments` | GET | ✗ | ✗ | ✓ | ✗ | ✓ | |

---

## assistant-svc (порт 8083)

| Эндпоинт | Метод | public | borrower | employee | expert | admin | Примечание |
|----------|-------|--------|---------|---------|--------|-------|-----------|
| `/assistant/chat` | POST | ✓ | ✓ | ✓ | ✓ | ✓ | Rate limit 20 req/min по IP |
| `/assistant/chat/ticket` | POST | ✗ | ✓ | ✗ | ✗ | ✗ | Только заёмщик создаёт тикет |

---

## Реализация в Go (middleware)

```go
// Пример RequireRole middleware для chi router
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := jwt.ClaimsFromContext(r.Context())
            for _, role := range roles {
                if claims.Role == role {
                    next.ServeHTTP(w, r)
                    return
                }
            }
            http.Error(w, `{"code":"FORBIDDEN"}`, http.StatusForbidden)
        })
    }
}

// Использование:
r.Get("/api/v1/programs", RequireRole("public", "borrower", "employee", "expert", "admin")(programsHandler))
r.Post("/api/v1/applications", RequireRole("borrower", "admin")(createApplicationHandler))
r.Get("/api/v1/admin/users", RequireRole("admin")(adminUsersHandler))
```

---

## Итого по доступам

| | public | borrower | employee | expert | admin |
|--|--------|---------|---------|--------|-------|
| Публичный сайт | ✓ | ✓ | ✓ | ✓ | ✓ |
| Личный кабинет | ✗ | ✓ | ✓ | ✗ | ✓ |
| Подача заявки | ✗ | ✓ | ✗ | ✗ | ✓ |
| Экспертиза | ✗ | ✗ | ✓* | ✓* | ✓ |
| КК (голосование) | ✗ | ✗ | ✗ | ✓ | ✓ |
| Аналитика | ✗ | ✗ | ✓ | ✗ | ✓ |
| Sync данные | ✗ | ✗ | ✓ | ✓ | ✓ |
| Управление пользователями | ✗ | ✗ | ✗ | ✗ | ✓ |
| Виртуальный ассистент | ✓ | ✓ | ✓ | ✓ | ✓ |
