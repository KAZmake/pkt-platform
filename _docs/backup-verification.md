# Backup Verification Checklist — pg_dump → B2 → Restore

## Prerequisites

- Staging k8s cluster running (`kubectl get nodes`)
- `pkt-backup-secrets` created: `B2_ACCOUNT_ID`, `B2_APPLICATION_KEY`, `B2_BUCKET`
- `rclone` installed locally for manual checks

---

## Part 1 — Trigger a manual backup

```bash
# Trigger CronJob manually
kubectl create job --from=cronjob/pg-backup pg-backup-manual -n staging

# Watch logs
kubectl logs -f job/pg-backup-manual -n staging

# Expected output:
# [2026-...] Starting backup — timestamp: 20260524_030000
#   Dumping pkt_db ...
#   pkt_db: 4.2M
#   Dumping keycloak_db ...
#   ...
#   Upload complete.
# [2026-...] Backup finished successfully.
```

## Part 2 — Verify files in B2

```bash
# Using rclone (configure [b2] remote first: rclone config)
rclone ls b2:pkt-staging-backups/postgres/ --b2-account=<id> --b2-key=<key>

# Expected: one directory per backup timestamp, 4 .dump.gz files each
# postgres/20260524_030000/pkt_db_20260524_030000.dump.gz
# postgres/20260524_030000/keycloak_db_20260524_030000.dump.gz
# postgres/20260524_030000/directus_db_20260524_030000.dump.gz
# postgres/20260524_030000/expertise_db_20260524_030000.dump.gz
```

## Part 3 — Restore to a test database

```bash
# Run restore in a temporary pod (NEVER against production)
kubectl run restore-test \
  --image=ghcr.io/kazmaker/pkt-platform/pg-backup:latest \
  --restart=Never \
  --rm -it \
  --env="PGHOST=postgres" \
  --env="PGUSER=pkt" \
  --env="PGPASSWORD=<from secret>" \
  --env="B2_BUCKET=pkt-staging-backups" \
  --env="B2_ACCOUNT_ID=<id>" \
  --env="B2_APPLICATION_KEY=<key>" \
  -n staging \
  -- restore 20260524_030000

# Script will prompt: "Type 'yes' to continue"
# After restore, verify row counts:
kubectl exec -n staging deploy/postgres -- \
  psql -U pkt -d pkt_db -c "SELECT COUNT(*) FROM users;"
```

## Part 4 — Verify row counts post-restore

```sql
-- Run inside the postgres pod
\c pkt_db
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM applications;
SELECT COUNT(*) FROM borrowers;

\c directus_db
SELECT COUNT(*) FROM directus_users;

\c expertise_db
-- tables vary, just check schema exists:
\dt
```

## Part 5 — Retention policy

After 30 days, old backups should be pruned automatically. Verify manually:

```bash
# List all backup dirs and confirm none older than 30 days
rclone lsd b2:pkt-staging-backups/postgres/ | awk '{print $NF}' | sort
```

## Part 6 — CronJob schedule check

```bash
kubectl get cronjob pg-backup -n staging
# Expected: SCHEDULE=0 3 * * *, SUSPEND=False, ACTIVE=0, LAST SCHEDULE=...

kubectl get jobs -n staging | grep pg-backup
# Most recent: COMPLETIONS=1/1
```

---

## Checklist

| #   | Check                                            | Result |
| --- | ------------------------------------------------ | ------ |
| 1   | Manual backup job completes without errors       | `[ ]`  |
| 2   | 4 dump files appear in B2 bucket                 | `[ ]`  |
| 3   | Each dump file is > 0 bytes                      | `[ ]`  |
| 4   | Restore completes without errors                 | `[ ]`  |
| 5   | Row counts match pre-restore values              | `[ ]`  |
| 6   | Retention pruning removes 31-day-old test backup | `[ ]`  |
| 7   | CronJob shows last successful run in kubectl     | `[ ]`  |
