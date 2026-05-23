#!/usr/bin/env python3
"""
Создание полей коллекций Directus для PKT Platform.
Запускается после создания коллекций через make directus-init.
"""
import json
import urllib.request
import urllib.error
import os
import sys

BASE = os.environ.get("DIRECTUS_URL", "http://localhost:8055")
EMAIL = os.environ.get("DIRECTUS_ADMIN_EMAIL", "admin@pkt.kz")
PASSWORD = os.environ.get("DIRECTUS_ADMIN_PASSWORD", "admin123")


def api(method, path, data=None, token=None):
    url = f"{BASE}{path}"
    body = json.dumps(data).encode() if data else None
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req) as r:
            return json.loads(r.read())
    except urllib.error.HTTPError as e:
        return json.loads(e.read())


# Auth
resp = api("POST", "/auth/login", {"email": EMAIL, "password": PASSWORD})
token = resp["data"]["access_token"]
print(f"Authenticated as {EMAIL}")


def add_field(collection, field, field_type, meta=None, schema=None):
    payload = {"field": field, "type": field_type}
    if meta:
        payload["meta"] = meta
    if schema:
        payload["schema"] = schema
    resp = api("POST", f"/fields/{collection}", payload, token)
    if "errors" in resp:
        err = resp["errors"][0]["message"] if resp["errors"] else "unknown"
        if "already exists" in err.lower():
            print(f"  {field}: exists")
        else:
            print(f"  {field}: ERROR — {err}")
    else:
        print(f"  {field}: OK")


# ── programs ──────────────────────────────────────────────────────────────────
print("\n--- programs fields ---")
add_field("programs", "name",            "string",  {"required": True, "note": "Название RU"},         {"is_nullable": False})
add_field("programs", "name_kz",         "string",  {"note": "Название KZ"})
add_field("programs", "name_en",         "string",  {"note": "Название EN"})
add_field("programs", "slug",            "string",  {"note": "URL slug"},                               {"is_unique": True})
add_field("programs", "description",     "text",    {"interface": "input-rich-text-md", "note": "Описание RU"})
add_field("programs", "description_kz",  "text",    {"interface": "input-rich-text-md", "note": "Описание KZ"})
add_field("programs", "rate",            "decimal", {"note": "% годовых"},                              {"numeric_precision": 5,  "numeric_scale": 2})
add_field("programs", "min_amount",      "decimal", {"note": "Минимальная сумма"},                      {"numeric_precision": 15, "numeric_scale": 2})
add_field("programs", "max_amount",      "decimal", {"note": "Максимальная сумма"},                     {"numeric_precision": 15, "numeric_scale": 2})
add_field("programs", "min_term_months", "integer", {"note": "Мин. срок (мес.)"})
add_field("programs", "max_term_months", "integer", {"note": "Макс. срок (мес.)"})
add_field("programs", "activity_types",  "json",    {"note": "crop_farming | livestock | mixed"})
add_field("programs", "is_active",       "boolean", {"note": "Активна"},                                {"default_value": True})
add_field("programs", "sort",            "integer", {"hidden": True})
add_field("programs", "status",          "string",
    {"interface": "select-dropdown", "options": {"choices": [
        {"text": "Опубликовано", "value": "published"},
        {"text": "Черновик",     "value": "draft"},
    ]}},
    {"default_value": "published"})

# ── news ──────────────────────────────────────────────────────────────────────
print("\n--- news fields ---")
add_field("news", "title",        "string",   {"required": True, "note": "Заголовок RU"}, {"is_nullable": False})
add_field("news", "title_kz",     "string",   {"note": "Заголовок KZ"})
add_field("news", "title_en",     "string",   {"note": "Заголовок EN"})
add_field("news", "slug",         "string",   {"note": "URL slug"},                       {"is_unique": True})
add_field("news", "content",      "text",     {"interface": "input-rich-text-md", "note": "Содержание RU"})
add_field("news", "content_kz",   "text",     {"interface": "input-rich-text-md", "note": "Содержание KZ"})
add_field("news", "image",        "uuid",     {"interface": "file-image", "note": "Главное изображение"},
    {"foreign_key_table": "directus_files", "foreign_key_column": "id"})
add_field("news", "published_at", "dateTime", {"note": "Дата публикации"})
add_field("news", "sort",         "integer",  {"hidden": True})
add_field("news", "status",       "string",
    {"interface": "select-dropdown", "options": {"choices": [
        {"text": "Опубликовано", "value": "published"},
        {"text": "Черновик",     "value": "draft"},
        {"text": "Архив",        "value": "archived"},
    ]}},
    {"default_value": "draft"})

# ── faq ───────────────────────────────────────────────────────────────────────
print("\n--- faq fields ---")
add_field("faq", "question",    "string",  {"required": True, "note": "Вопрос RU"}, {"is_nullable": False})
add_field("faq", "question_kz", "string",  {"note": "Вопрос KZ"})
add_field("faq", "answer",      "text",    {"interface": "input-rich-text-md", "note": "Ответ RU"})
add_field("faq", "answer_kz",   "text",    {"interface": "input-rich-text-md", "note": "Ответ KZ"})
add_field("faq", "category",    "string",
    {"interface": "select-dropdown", "options": {"choices": [
        {"text": "Общие",      "value": "general"},
        {"text": "Займы",      "value": "loans"},
        {"text": "Документы",  "value": "documents"},
        {"text": "Прочее",     "value": "other"},
    ]}},
    {"default_value": "general"})
add_field("faq", "sort",   "integer", {"hidden": True})
add_field("faq", "status", "string",
    {"interface": "select-dropdown", "options": {"choices": [
        {"text": "Опубликовано", "value": "published"},
        {"text": "Черновик",     "value": "draft"},
    ]}},
    {"default_value": "published"})

# ── projects ──────────────────────────────────────────────────────────────────
print("\n--- projects fields ---")
add_field("projects", "title",         "string",  {"required": True, "note": "Название RU"}, {"is_nullable": False})
add_field("projects", "title_kz",      "string",  {"note": "Название KZ"})
add_field("projects", "description",   "text",    {"interface": "input-rich-text-md", "note": "Описание RU"})
add_field("projects", "description_kz","text",    {"interface": "input-rich-text-md", "note": "Описание KZ"})
add_field("projects", "image",         "uuid",    {"interface": "file-image", "note": "Фото проекта"},
    {"foreign_key_table": "directus_files", "foreign_key_column": "id"})
add_field("projects", "borrower_name", "string",  {"note": "Имя заёмщика"})
add_field("projects", "location",      "string",  {"note": "Район / населённый пункт"})
add_field("projects", "program_name",  "string",  {"note": "Программа кредитования"})
add_field("projects", "amount",        "decimal", {"note": "Сумма займа"},                   {"numeric_precision": 15, "numeric_scale": 2})
add_field("projects", "year",          "integer", {"note": "Год реализации"})
add_field("projects", "sort",          "integer", {"hidden": True})
add_field("projects", "status",        "string",
    {"interface": "select-dropdown", "options": {"choices": [
        {"text": "Опубликовано", "value": "published"},
        {"text": "Черновик",     "value": "draft"},
    ]}},
    {"default_value": "published"})

# ── Итог ──────────────────────────────────────────────────────────────────────
print("\n=== Done ===")
resp = api("GET", "/collections", token=token)
custom = [c["collection"] for c in resp["data"] if not c["collection"].startswith("directus_")]
print("Custom collections:", custom)
