.PHONY: up down restart logs ps health clean migrate migrate-down nats-init directus-init

COMPOSE  = docker compose -f infra/docker-compose.yml --env-file infra/.env.docker
DB_URL       = postgres://pkt:pkt_secret@localhost:5433/pkt_db?sslmode=disable
MIGRATE_CORE      = postgres://pkt:pkt_secret@pkt-postgres:5432/pkt_db?sslmode=disable&x-migrations-table=schema_migrations_core
MIGRATE_EXPERTISE = postgres://pkt:pkt_secret@pkt-postgres:5432/pkt_db?sslmode=disable&x-migrations-table=schema_migrations_expertise
MIGRATE_SYNC      = postgres://pkt:pkt_secret@pkt-postgres:5432/pkt_db?sslmode=disable&x-migrations-table=schema_migrations_sync
MIGRATE      = docker run --rm --network infra_pkt-net \
               -v $(PWD)/apps/api/core-api/migrations:/core \
               -v $(PWD)/apps/api/expertise-svc/migrations:/expertise \
               -v $(PWD)/apps/api/sync-svc/migrations:/sync \
               migrate/migrate

# ─── Основные команды ─────────────────────────────────────────────────────────

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) restart

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

# ─── Отдельные сервисы ────────────────────────────────────────────────────────

up-db:
	$(COMPOSE) up -d postgres valkey

up-infra:
	$(COMPOSE) up -d postgres valkey minio nats

logs-%:
	$(COMPOSE) logs -f $*

# ─── Health check всех сервисов ───────────────────────────────────────────────

health:
	@echo "=== PostgreSQL ===" && docker exec pkt-postgres pg_isready -U pkt -d pkt_db
	@echo "=== Valkey ===" && docker exec pkt-valkey valkey-cli ping
	@echo "=== MinIO ===" && docker exec pkt-minio mc ready local 2>/dev/null || echo "MinIO: check console at http://localhost:9001"
	@echo "=== NATS ===" && curl -s http://localhost:8222/healthz | grep -o '"status":"ok"' || echo "NATS: not ready"
	@echo "=== Keycloak ===" && curl -sf http://localhost:8180/realms/master > /dev/null && echo "OK" || echo "NOT READY"
	@echo "=== Directus ===" && curl -sf http://localhost:8055/server/health > /dev/null && echo "OK" || echo "NOT READY"

# ─── NATS JetStream стримы ────────────────────────────────────────────────────

nats-init:
	docker run --rm --network infra_pkt-net \
	  -v $(PWD)/infra/nats/streams-init.sh:/streams-init.sh:ro \
	  --entrypoint /bin/sh \
	  natsio/nats-box:latest /streams-init.sh

# ─── Directus коллекции ───────────────────────────────────────────────────────

directus-init:
	DIRECTUS_URL=http://localhost:8055 \
	DIRECTUS_ADMIN_EMAIL=$${DIRECTUS_ADMIN_EMAIL:-admin@pkt.kz} \
	DIRECTUS_ADMIN_PASSWORD=$${DIRECTUS_ADMIN_PASSWORD:-admin123} \
	python3 infra/directus/collections-init.py

# ─── Миграции (golang-migrate через Docker) ───────────────────────────────────
# Порядок важен: core-api (users, geo, borrowers) → expertise → sync

migrate:
	@echo "=== Migrating core-api ===" && \
	 $(MIGRATE) -path /core      -database "$(MIGRATE_CORE)"      up
	@echo "=== Migrating expertise-svc ===" && \
	 $(MIGRATE) -path /expertise -database "$(MIGRATE_EXPERTISE)" up
	@echo "=== Migrating sync-svc ===" && \
	 $(MIGRATE) -path /sync      -database "$(MIGRATE_SYNC)"      up
	@echo "All migrations applied."

migrate-down:
	@echo "=== Rolling back sync-svc ===" && \
	 $(MIGRATE) -path /sync      -database "$(MIGRATE_SYNC)"      down 1
	@echo "=== Rolling back expertise-svc ===" && \
	 $(MIGRATE) -path /expertise -database "$(MIGRATE_EXPERTISE)" down 1
	@echo "=== Rolling back core-api ===" && \
	 $(MIGRATE) -path /core      -database "$(MIGRATE_CORE)"      down 1

# ─── Очистка данных (ОСТОРОЖНО — удаляет volumes) ────────────────────────────

clean:
	$(COMPOSE) down -v
	@echo "All volumes removed."
