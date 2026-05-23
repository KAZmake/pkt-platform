#!/bin/sh
# Создание коллекций Directus для PKT Platform
# Запускается один раз: make directus-init

set -e

BASE="http://directus:8055"

echo "=== Authenticating with Directus ==="
TOKEN=$(wget -qO- --post-data='{"email":"'"$DIRECTUS_ADMIN_EMAIL"'","password":"'"$DIRECTUS_ADMIN_PASSWORD"'"}' \
  --header='Content-Type: application/json' \
  "$BASE/auth/login" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")

if [ -z "$TOKEN" ]; then
  echo "ERROR: Could not authenticate with Directus"
  exit 1
fi

echo "Authenticated OK."

create_collection() {
  COLLECTION=$1
  PAYLOAD=$2
  STATUS=$(wget -qO- --method=POST \
    --header="Authorization: Bearer $TOKEN" \
    --header='Content-Type: application/json' \
    --body-data="$PAYLOAD" \
    --server-response \
    "$BASE/collections" 2>&1 | grep "HTTP/" | tail -1 | awk '{print $2}')
  if [ "$STATUS" = "200" ] || [ "$STATUS" = "204" ]; then
    echo "Collection '$COLLECTION': created."
  else
    echo "Collection '$COLLECTION': already exists or error ($STATUS), skipping."
  fi
}

create_field() {
  COLLECTION=$1
  PAYLOAD=$2
  wget -qO- --method=POST \
    --header="Authorization: Bearer $TOKEN" \
    --header='Content-Type: application/json' \
    --body-data="$PAYLOAD" \
    "$BASE/fields/$COLLECTION" > /dev/null 2>&1 || true
}

# ─── programs (Программы кредитования) ───────────────────────────────────────
echo "--- Creating collection: programs ---"
create_collection "programs" '{
  "collection": "programs",
  "meta": {
    "icon": "payments",
    "note": "Программы кредитования — публичный справочник",
    "display_template": "{{name}}",
    "sort_field": "sort"
  },
  "schema": {}
}'

for FIELD in \
  '{"field":"name","type":"string","meta":{"required":true,"note":"Название RU"},"schema":{"is_nullable":false}}' \
  '{"field":"name_kz","type":"string","meta":{"note":"Название KZ"}}' \
  '{"field":"name_en","type":"string","meta":{"note":"Название EN"}}' \
  '{"field":"slug","type":"string","meta":{"note":"URL slug"},"schema":{"is_unique":true}}' \
  '{"field":"description","type":"text","meta":{"interface":"input-rich-text-md","note":"Описание RU"}}' \
  '{"field":"description_kz","type":"text","meta":{"interface":"input-rich-text-md","note":"Описание KZ"}}' \
  '{"field":"rate","type":"decimal","meta":{"note":"% годовых"},"schema":{"numeric_precision":5,"numeric_scale":2}}' \
  '{"field":"min_amount","type":"decimal","meta":{"note":"Минимальная сумма"},"schema":{"numeric_precision":15,"numeric_scale":2}}' \
  '{"field":"max_amount","type":"decimal","meta":{"note":"Максимальная сумма"},"schema":{"numeric_precision":15,"numeric_scale":2}}' \
  '{"field":"min_term_months","type":"integer","meta":{"note":"Мин. срок (мес.)"}}' \
  '{"field":"max_term_months","type":"integer","meta":{"note":"Макс. срок (мес.)"}}' \
  '{"field":"activity_types","type":"json","meta":{"note":"Типы деятельности: crop_farming|livestock|mixed"}}' \
  '{"field":"is_active","type":"boolean","meta":{"note":"Активна"},"schema":{"default_value":true}}' \
  '{"field":"sort","type":"integer","meta":{"hidden":true}}' \
  '{"field":"status","type":"string","meta":{"interface":"select-dropdown","options":{"choices":[{"text":"Опубликовано","value":"published"},{"text":"Черновик","value":"draft"}]}},"schema":{"default_value":"published"}}'
do
  create_field "programs" "$FIELD"
done
echo "Collection 'programs': fields created."

# ─── news (Новости) ───────────────────────────────────────────────────────────
echo "--- Creating collection: news ---"
create_collection "news" '{
  "collection": "news",
  "meta": {
    "icon": "article",
    "note": "Новости компании",
    "display_template": "{{title}}",
    "sort_field": "sort"
  },
  "schema": {}
}'

for FIELD in \
  '{"field":"title","type":"string","meta":{"required":true,"note":"Заголовок RU"},"schema":{"is_nullable":false}}' \
  '{"field":"title_kz","type":"string","meta":{"note":"Заголовок KZ"}}' \
  '{"field":"title_en","type":"string","meta":{"note":"Заголовок EN"}}' \
  '{"field":"slug","type":"string","meta":{"note":"URL slug"},"schema":{"is_unique":true}}' \
  '{"field":"content","type":"text","meta":{"interface":"input-rich-text-md","note":"Содержание RU"}}' \
  '{"field":"content_kz","type":"text","meta":{"interface":"input-rich-text-md","note":"Содержание KZ"}}' \
  '{"field":"image","type":"uuid","meta":{"interface":"file-image","note":"Главное изображение"},"schema":{"foreign_key_table":"directus_files","foreign_key_column":"id"}}' \
  '{"field":"published_at","type":"dateTime","meta":{"note":"Дата публикации"}}' \
  '{"field":"sort","type":"integer","meta":{"hidden":true}}' \
  '{"field":"status","type":"string","meta":{"interface":"select-dropdown","options":{"choices":[{"text":"Опубликовано","value":"published"},{"text":"Черновик","value":"draft"},{"text":"Архив","value":"archived"}]}},"schema":{"default_value":"draft"}}'
do
  create_field "news" "$FIELD"
done
echo "Collection 'news': fields created."

# ─── faq (Часто задаваемые вопросы) ──────────────────────────────────────────
echo "--- Creating collection: faq ---"
create_collection "faq" '{
  "collection": "faq",
  "meta": {
    "icon": "help",
    "note": "Часто задаваемые вопросы",
    "display_template": "{{question}}",
    "sort_field": "sort"
  },
  "schema": {}
}'

for FIELD in \
  '{"field":"question","type":"string","meta":{"required":true,"note":"Вопрос RU"},"schema":{"is_nullable":false}}' \
  '{"field":"question_kz","type":"string","meta":{"note":"Вопрос KZ"}}' \
  '{"field":"answer","type":"text","meta":{"interface":"input-rich-text-md","note":"Ответ RU"}}' \
  '{"field":"answer_kz","type":"text","meta":{"interface":"input-rich-text-md","note":"Ответ KZ"}}' \
  '{"field":"category","type":"string","meta":{"interface":"select-dropdown","options":{"choices":[{"text":"Общие","value":"general"},{"text":"Займы","value":"loans"},{"text":"Документы","value":"documents"},{"text":"Прочее","value":"other"}]}},"schema":{"default_value":"general"}}' \
  '{"field":"sort","type":"integer","meta":{"hidden":true}}' \
  '{"field":"status","type":"string","meta":{"interface":"select-dropdown","options":{"choices":[{"text":"Опубликовано","value":"published"},{"text":"Черновик","value":"draft"}]}},"schema":{"default_value":"published"}}'
do
  create_field "faq" "$FIELD"
done
echo "Collection 'faq': fields created."

# ─── projects (Реализованные проекты) ────────────────────────────────────────
echo "--- Creating collection: projects ---"
create_collection "projects" '{
  "collection": "projects",
  "meta": {
    "icon": "agriculture",
    "note": "Реализованные проекты заёмщиков",
    "display_template": "{{title}}",
    "sort_field": "sort"
  },
  "schema": {}
}'

for FIELD in \
  '{"field":"title","type":"string","meta":{"required":true,"note":"Название RU"},"schema":{"is_nullable":false}}' \
  '{"field":"title_kz","type":"string","meta":{"note":"Название KZ"}}' \
  '{"field":"description","type":"text","meta":{"interface":"input-rich-text-md","note":"Описание RU"}}' \
  '{"field":"description_kz","type":"text","meta":{"interface":"input-rich-text-md","note":"Описание KZ"}}' \
  '{"field":"image","type":"uuid","meta":{"interface":"file-image","note":"Фото проекта"},"schema":{"foreign_key_table":"directus_files","foreign_key_column":"id"}}' \
  '{"field":"borrower_name","type":"string","meta":{"note":"Имя заёмщика"}}' \
  '{"field":"location","type":"string","meta":{"note":"Район / населённый пункт"}}' \
  '{"field":"program_name","type":"string","meta":{"note":"Программа кредитования"}}' \
  '{"field":"amount","type":"decimal","meta":{"note":"Сумма займа"},"schema":{"numeric_precision":15,"numeric_scale":2}}' \
  '{"field":"year","type":"integer","meta":{"note":"Год реализации"}}' \
  '{"field":"sort","type":"integer","meta":{"hidden":true}}' \
  '{"field":"status","type":"string","meta":{"interface":"select-dropdown","options":{"choices":[{"text":"Опубликовано","value":"published"},{"text":"Черновик","value":"draft"}]}},"schema":{"default_value":"published"}}'
do
  create_field "projects" "$FIELD"
done
echo "Collection 'projects': fields created."

echo ""
echo "=== Directus collections initialized ==="
wget -qO- --header="Authorization: Bearer $TOKEN" "$BASE/collections" \
  | python3 -c "
import sys,json
cols = [c['collection'] for c in json.load(sys.stdin)['data'] if not c['collection'].startswith('directus_')]
print('Custom collections:', cols)
"
