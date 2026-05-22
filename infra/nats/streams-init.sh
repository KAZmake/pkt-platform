#!/bin/sh
# Создание JetStream стримов для PKT Platform
# Запускается один раз: make nats-init

set -e

SERVER="nats://nats:4222"

echo "=== Creating NATS JetStream streams ==="

# ─── application-events ───────────────────────────────────────────────────────
# Публикует: expertise-svc при каждом переходе FSM заявки
# Консьюмеры: core-api (уведомления), sync-svc (после выдачи займа)
cat > /tmp/stream-application-events.json <<'EOF'
{
  "name": "application-events",
  "subjects": ["application.>"],
  "storage": "file",
  "retention": "limits",
  "max_age": 2592000000000000,
  "max_bytes": 1073741824,
  "max_msg_size": 1048576,
  "max_msgs": -1,
  "num_replicas": 1,
  "discard": "old",
  "duplicate_window": 120000000000
}
EOF

if nats stream info application-events --server "$SERVER" > /dev/null 2>&1; then
  echo "Stream 'application-events': already exists, skipping."
else
  nats stream add --server "$SERVER" --config /tmp/stream-application-events.json
  echo "Stream 'application-events': created."
fi

# ─── notifications ────────────────────────────────────────────────────────────
# Публикует: core-api (статус, платежи, тикеты, системные)
# Консьюмеры: notification worker → Resend email + WebSocket push
cat > /tmp/stream-notifications.json <<'EOF'
{
  "name": "notifications",
  "subjects": ["notifications.>"],
  "storage": "file",
  "retention": "limits",
  "max_age": 604800000000000,
  "max_bytes": 536870912,
  "max_msg_size": 262144,
  "max_msgs": 1000000,
  "num_replicas": 1,
  "discard": "old",
  "duplicate_window": 120000000000
}
EOF

if nats stream info notifications --server "$SERVER" > /dev/null 2>&1; then
  echo "Stream 'notifications': already exists, skipping."
else
  nats stream add --server "$SERVER" --config /tmp/stream-notifications.json
  echo "Stream 'notifications': created."
fi

# ─── sync-events ──────────────────────────────────────────────────────────────
# Публикует: sync-svc после синхронизации с 1С (каждые 20 мин)
# Консьюмеры: core-api (обновление кэша займов), expertise-svc (алерты залогов)
cat > /tmp/stream-sync-events.json <<'EOF'
{
  "name": "sync-events",
  "subjects": ["sync.>"],
  "storage": "file",
  "retention": "limits",
  "max_age": 172800000000000,
  "max_bytes": 268435456,
  "max_msg_size": 524288,
  "max_msgs": 100000,
  "num_replicas": 1,
  "discard": "old",
  "duplicate_window": 120000000000
}
EOF

if nats stream info sync-events --server "$SERVER" > /dev/null 2>&1; then
  echo "Stream 'sync-events': already exists, skipping."
else
  nats stream add --server "$SERVER" --config /tmp/stream-sync-events.json
  echo "Stream 'sync-events': created."
fi

echo ""
echo "=== NATS streams initialized ==="
nats stream ls --server "$SERVER"
