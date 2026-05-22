#!/bin/sh
# Задача 1.4: создание бакетов MinIO при первом запуске

set -e

MC="mc"
ALIAS="pkt"
ENDPOINT="http://minio:9000"

echo "Waiting for MinIO to be ready..."
until $MC alias set $ALIAS $ENDPOINT "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD" 2>/dev/null; do
  sleep 2
done
echo "MinIO is ready."

# ─── Создаём бакеты ──────────────────────────────────────────────────────────
for BUCKET in documents uploads media templates; do
  if $MC ls "$ALIAS/$BUCKET" > /dev/null 2>&1; then
    echo "Bucket '$BUCKET' already exists, skipping."
  else
    $MC mb "$ALIAS/$BUCKET"
    echo "Bucket '$BUCKET' created."
  fi
done

# ─── Политики доступа ────────────────────────────────────────────────────────

# media — публичное чтение (публичные изображения сайта)
$MC anonymous set download "$ALIAS/media"
echo "Bucket 'media' set to public download."

# documents, uploads, templates — приватные (только через presigned URL)
$MC anonymous set none "$ALIAS/documents"
$MC anonymous set none "$ALIAS/uploads"
$MC anonymous set none "$ALIAS/templates"
echo "Buckets 'documents', 'uploads', 'templates' set to private."

echo "MinIO init complete."
