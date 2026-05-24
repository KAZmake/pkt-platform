#!/usr/bin/env bash
# Restore PostgreSQL databases from a B2 backup.
#
# Usage:
#   ./restore.sh <backup-timestamp>
#   ./restore.sh 20260524_030000
#
# Environment variables: same as backup.sh
# WARNING: This will DROP and recreate the target databases.

set -euo pipefail

TIMESTAMP="${1:?Usage: $0 <backup-timestamp> (e.g. 20260524_030000)}"
B2_BUCKET="${B2_BUCKET:?B2_BUCKET is required}"
BACKUP_PREFIX="postgres"
DATABASES=(pkt_db keycloak_db directus_db expertise_db)

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "[$(date -u)] Starting restore from timestamp: ${TIMESTAMP}"
echo "  WARNING: target databases will be dropped and recreated!"
echo ""
read -rp "  Type 'yes' to continue: " CONFIRM
[[ "$CONFIRM" == "yes" ]] || { echo "Aborted."; exit 1; }

# ── Download from B2 ───────────────────────────────────────────────────────
echo "  Downloading from b2:${B2_BUCKET}/${BACKUP_PREFIX}/${TIMESTAMP}/ ..."

rclone copy "b2:${B2_BUCKET}/${BACKUP_PREFIX}/${TIMESTAMP}/" "${TMPDIR}/" \
  --b2-account="${B2_ACCOUNT_ID}" \
  --b2-key="${B2_APPLICATION_KEY}"

echo "  Downloaded files:"
ls -lh "${TMPDIR}/"

# ── Restore each database ─────────────────────────────────────────────────
for DB in "${DATABASES[@]}"; do
  DUMP_FILE="${TMPDIR}/${DB}_${TIMESTAMP}.dump.gz"

  if [[ ! -f "${DUMP_FILE}" ]]; then
    echo "  SKIP: ${DUMP_FILE} not found"
    continue
  fi

  echo "  Restoring ${DB} ..."

  # Drop and recreate
  PGPASSWORD="${PGPASSWORD}" psql \
    --host="${PGHOST:-postgres}" \
    --port="${PGPORT:-5432}" \
    --username="${PGUSER:-pkt}" \
    --dbname=postgres \
    --no-password \
    -c "DROP DATABASE IF EXISTS ${DB}; CREATE DATABASE ${DB};" 2>&1

  # Restore
  zcat "${DUMP_FILE}" \
    | PGPASSWORD="${PGPASSWORD}" pg_restore \
        --host="${PGHOST:-postgres}" \
        --port="${PGPORT:-5432}" \
        --username="${PGUSER:-pkt}" \
        --dbname="${DB}" \
        --no-password \
        --no-owner \
        --format=custom \
        --exit-on-error 2>&1

  echo "  ${DB}: restored OK"
done

echo ""
echo "[$(date -u)] Restore complete."
