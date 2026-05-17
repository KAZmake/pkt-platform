# ТОО «Первое кредитное товарищество» — Веб-платформа + Мобильное приложение

## ЧИТАЙ ПЕРВЫМ ПРИ КАЖДОЙ СЕССИИ

Это финансовая платформа для кредитного товарищества (агросектор, ЗКО, Казахстан).
Подробности — в `_docs/architecture.md`, `_docs/roadmap.md`, `_docs/decisions.md`.

---

## Монорепозиторий (Turborepo)

```
apps/
  web/          → Next.js 14+ (сайт)
  mobile/       → React Native + Expo (iOS + Android)
  api/          → Go микросервисы
packages/
  shared/       → типы, API-клиент, утилиты (переиспользуется web + mobile)
_docs/          → контекст проекта (этот файл и остальные)
```

---

## Стек — кратко

| Слой | Технология |
|------|-----------|
| Frontend | Next.js 14+ App Router + TypeScript + TailwindCSS |
| Mobile | React Native + Expo (NOT WebView) |
| Backend | Go + Fiber/chi + gRPC между сервисами |
| Queue | NATS JetStream |
| DB | PostgreSQL + TimescaleDB + PostGIS |
| Cache | Redis (Valkey) |
| Files | MinIO (S3-совместимый, on-premise) |
| Auth | Keycloak (RBAC + OIDC + SAML 2.0 + MFA) |
| CMS | Directus (headless, REST+GraphQL) |
| Analytics | Grafana + Prometheus + Loki |
| Infra | Oracle Cloud Always Free → Docker + k3s → Kubernetes |
| CDN/DDoS | Cloudflare |
| Email (users) | Resend + React Email |
| Email (staff) | Microsoft 365 или Google Workspace (SAML→Keycloak) |
| Map | MapLibre GL JS / Mapbox GL JS |
| Backups | Backblaze B2 + pg_dump cron |
| Ext. monitor | UptimeRobot |
| AI assistant | Anthropic API (Claude) |

API версионирование: `/api/v1/...` — обязательно.

---

## Роли пользователей

| Роль | Кодовое имя | Доступ |
|------|-------------|--------|
| Публичный посетитель | `public` | Всё публичное: сайт, карта (публ. слои), калькулятор, чат-бот |
| Заёмщик | `borrower` | + Личный кабинет (документы, график, обращения, уведомления) |
| Сотрудник | `employee` | + Аналитика, Экспертиза (свой этап), карта (все слои), инлайн-редактор, Почта и диск |
| Эксперт / КК | `expert` | + экспертиза залогов, юр. проверка, протоколы КК, голосование |
| Администратор | `admin` | Полный доступ, управление пользователями, системные логи |

Keycloak: MFA включён, SAML 2.0 для Microsoft 365/Google Workspace, PKCE для мобилки.
Каждый API-эндпоинт защищён middleware проверки JWT + роли.

---

## Навигация сайта

**Публичная:** Главная · О компании · Программы кредитования · Карта · Проекты · Новости · Контакты · FAQ

**Хедер (всегда):** Логотип · Номер Call-центра · Поиск · Версия для слабовидящих · Смена языка (Қаз/Рус/Eng) · Личный кабинет · Подать заявку · Чат-бот (виджет правый нижний угол)

**После авторизации сотрудника:** + Аналитика · Экспертиза · Почта и диск

---

## Go-сервисы (отдельные микросервисы)

1. **core-api** — пользователи, программы, заявки, ЛК, уведомления
2. **expertise-svc** — BPM-модуль экспертизы (FSM заявок)
3. **sync-svc** — синхронизация с 1С (cron каждые 20 мин, кэш Redis)
4. **assistant-svc** — прокси к Anthropic API (rate limit по IP)

---

## Критичные бизнес-правила (не нарушать)

- Финансовые данные — только из 1С через sync-svc, не хранить дубли без кэша
- Аудит-лог Экспертизы: только INSERT, никогда UPDATE/DELETE
- Карточка залога — самостоятельная сущность (высвобождение, привязка к новым кредитам)
- Карточка заёмщика — источник для формирования заявки
- Документы подписываются онлайн или прикрепляется скан
- После выдачи кредита — передача данных в 1С и в БД

---

## Текущий прогресс → см. `_docs/roadmap.md`
