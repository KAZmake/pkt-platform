#!/usr/bin/env bash
# pg_dump → compress → upload to Backblaze B2
#
# Environment variables required:
#   PGHOST, PGPORT, PGUSER, PGPASSWORD
#   B2_ACCOUNT_ID, B2_APPLICATION_KEY, B2_BUCKET
#   RETENTION_DAYS (default: 30)
#
# Dependencies: pg_dump, gzip, rclone (configured with [b2] remote)
# In staging k8s: runs as a CronJob with pkt-backup image.

set -euo pipefail

TIMESTAMP="$(date -u +%Y%m%d_%H%M%S)"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
B2_BUCKET="${B2_BUCKET:?B2_BUCKET is required}"
BACKUP_PREFIX="postgres"

# Databases to back up
DATABASES=(pkt_db keycloak_db directus_db expertise_db)

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "[$(date -u)] Starting backup — timestamp: ${TIMESTAMP}"

# ── Dump each database ─────────────────────────────────────────────────────
for DB in "${DATABASES[@]}"; do
  DUMP_FILE="${TMPDIR}/${DB}_${TIMESTAMP}.dump.gz"
  echo "  Dumping ${DB} ..."

  PGPASSWORD="${PGPASSWORD}" pg_dump \
    --host="${PGHOST:-postgres}" \
    --port="${PGPORT:-5432}" \
    --username="${PGUSER:-pkt}" \
    --format=custom \
    --no-password \
    "${DB}" \
    | gzip > "${DUMP_FILE}"

  SIZE=$(du -sh "${DUMP_FILE}" | cut -f1)
  echo "  ${DB}: ${SIZE}"
done

# ── Upload to B2 via rclone ────────────────────────────────────────────────
echo "  Uploading to b2:${B2_BUCKET}/${BACKUP_PREFIX}/${TIMESTAMP}/ ..."

rclone copy "${TMPDIR}/" "b2:${B2_BUCKET}/${BACKUP_PREFIX}/${TIMESTAMP}/" \
  --b2-account="${B2_ACCOUNT_ID}" \
  --b2-key="${B2_APPLICATION_KEY}" \
  --transfers=4 \
  --stats=10s

echo "  Upload complete."

# ── Prune old backups ──────────────────────────────────────────────────────
echo "  Pruning backups older than ${RETENTION_DAYS} days ..."

# List all backup dirs, delete those older than retention window
CUTOFF=$(date -u -d "${RETENTION_DAYS} days ago" +%Y%m%d 2>/dev/null \
  || date -u -v -"${RETENTION_DAYS}"d +%Y%m%d)  # macOS fallback

rclone lsd "b2:${B2_BUCKET}/${BACKUP_PREFIX}/" \
  --b2-account="${B2_ACCOUNT_ID}" \
  --b2-key="${B2_APPLICATION_KEY}" \
  | awk '{print $NF}' \
  | while read -r DIR; do
      DIR_DATE="${DIR%%_*}"  # extract YYYYMMDD prefix
      if [[ "${DIR_DATE}" < "${CUTOFF}" ]]; then
        echo "    Deleting old backup: ${DIR}"
        rclone purge "b2:${B2_BUCKET}/${BACKUP_PREFIX}/${DIR}/" \
          --b2-account="${B2_ACCOUNT_ID}" \
          --b2-key="${B2_APPLICATION_KEY}"
      fi
    done

echo "[$(date -u)] Backup finished successfully."
