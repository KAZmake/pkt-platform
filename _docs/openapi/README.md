# OpenAPI 3.0 — PKT Platform

Спецификации разбиты по сервисам. Все эндпоинты — `/api/v1/...`.

| Файл | Сервис | Порт | Описание |
|------|--------|------|---------|
| [core-api.yaml](./core-api.yaml) | core-api | 8080 | Пользователи, программы, заявки, ЛК, уведомления |
| [expertise-svc.yaml](./expertise-svc.yaml) | expertise-svc | 8081 | BPM экспертизы, залоги, КК |
| [sync-svc.yaml](./sync-svc.yaml) | sync-svc | 8082 | Синхронизация с 1С, мониторинг |
| [assistant-svc.yaml](./assistant-svc.yaml) | assistant-svc | 8083 | Виртуальный ассистент (Claude API) |

## Просмотр в Swagger UI

Вставь содержимое любого файла на [editor.swagger.io](https://editor.swagger.io).
