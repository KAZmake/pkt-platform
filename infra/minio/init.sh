#!/bin/sh
# MinIO init: создание бакетов + IAM-политики по ролям

set -e

MC="mc"
ALIAS="pkt"
ENDPOINT="http://minio:9000"

echo "Waiting for MinIO to be ready..."
until $MC alias set $ALIAS $ENDPOINT "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD" 2>/dev/null; do
  sleep 2
done
echo "MinIO is ready."

# ─── Бакеты ──────────────────────────────────────────────────────────────────

for BUCKET in documents uploads media templates; do
  if $MC ls "$ALIAS/$BUCKET" > /dev/null 2>&1; then
    echo "Bucket '$BUCKET' already exists, skipping."
  else
    $MC mb "$ALIAS/$BUCKET"
    echo "Bucket '$BUCKET' created."
  fi
done

# media — публичное чтение (изображения, публичные файлы сайта)
$MC anonymous set download "$ALIAS/media"
echo "Bucket 'media': public download."

# documents, uploads, templates — приватные (доступ только через presigned URL или IAM)
$MC anonymous set none "$ALIAS/documents"
$MC anonymous set none "$ALIAS/uploads"
$MC anonymous set none "$ALIAS/templates"
echo "Buckets documents/uploads/templates: private."

# ─── IAM политики по ролям ───────────────────────────────────────────────────

# borrower: загрузка в uploads/, чтение своих документов через presigned URL (выдаёт core-api)
cat > /tmp/policy-borrower.json <<'POLICY'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject"],
      "Resource": ["arn:aws:s3:::uploads/*"]
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject"],
      "Resource": [
        "arn:aws:s3:::documents/*",
        "arn:aws:s3:::media/*"
      ]
    }
  ]
}
POLICY
$MC admin policy create $ALIAS pkt-borrower /tmp/policy-borrower.json 2>/dev/null || \
$MC admin policy update $ALIAS pkt-borrower /tmp/policy-borrower.json 2>/dev/null || true
echo "Policy 'pkt-borrower' applied."

# employee: чтение всего, запись в documents/ и uploads/, чтение templates/
cat > /tmp/policy-employee.json <<'POLICY'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:ListBucket"],
      "Resource": [
        "arn:aws:s3:::documents",
        "arn:aws:s3:::documents/*",
        "arn:aws:s3:::uploads",
        "arn:aws:s3:::uploads/*",
        "arn:aws:s3:::media",
        "arn:aws:s3:::media/*",
        "arn:aws:s3:::templates",
        "arn:aws:s3:::templates/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:DeleteObject"],
      "Resource": [
        "arn:aws:s3:::documents/*",
        "arn:aws:s3:::uploads/*",
        "arn:aws:s3:::media/*"
      ]
    }
  ]
}
POLICY
$MC admin policy create $ALIAS pkt-employee /tmp/policy-employee.json 2>/dev/null || \
$MC admin policy update $ALIAS pkt-employee /tmp/policy-employee.json 2>/dev/null || true
echo "Policy 'pkt-employee' applied."

# expert: чтение documents/ и uploads/, запись заключений в documents/conclusions/
cat > /tmp/policy-expert.json <<'POLICY'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:ListBucket"],
      "Resource": [
        "arn:aws:s3:::documents",
        "arn:aws:s3:::documents/*",
        "arn:aws:s3:::uploads",
        "arn:aws:s3:::uploads/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject"],
      "Resource": ["arn:aws:s3:::documents/conclusions/*"]
    }
  ]
}
POLICY
$MC admin policy create $ALIAS pkt-expert /tmp/policy-expert.json 2>/dev/null || \
$MC admin policy update $ALIAS pkt-expert /tmp/policy-expert.json 2>/dev/null || true
echo "Policy 'pkt-expert' applied."

# admin: полный доступ ко всем бакетам
cat > /tmp/policy-admin.json <<'POLICY'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:*"],
      "Resource": [
        "arn:aws:s3:::*"
      ]
    }
  ]
}
POLICY
$MC admin policy create $ALIAS pkt-admin /tmp/policy-admin.json 2>/dev/null || \
$MC admin policy update $ALIAS pkt-admin /tmp/policy-admin.json 2>/dev/null || true
echo "Policy 'pkt-admin' applied."

# ─── Сервисные аккаунты для Go-сервисов ──────────────────────────────────────
# core-api: чтение/запись documents/ и uploads/ (presigned URL generation)
# expertise-svc: чтение documents/, запись conclusions/

$MC admin user add $ALIAS svc-core-api svc_core_api_secret 2>/dev/null || true
$MC admin policy attach $ALIAS pkt-employee --user svc-core-api 2>/dev/null || true
echo "Service account 'svc-core-api' ready."

$MC admin user add $ALIAS svc-expertise svc_expertise_secret 2>/dev/null || true
$MC admin policy attach $ALIAS pkt-expert --user svc-expertise 2>/dev/null || true
echo "Service account 'svc-expertise' ready."

echo ""
echo "MinIO init complete."
