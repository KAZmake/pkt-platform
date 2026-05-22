.PHONY: up down restart logs ps health clean

COMPOSE = docker compose -f infra/docker-compose.yml --env-file infra/.env.docker

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
	@echo "=== Keycloak ===" && curl -sf http://localhost:8180/health/ready > /dev/null && echo "OK" || echo "NOT READY"
	@echo "=== Directus ===" && curl -sf http://localhost:8055/server/health > /dev/null && echo "OK" || echo "NOT READY"

# ─── Очистка данных (ОСТОРОЖНО — удаляет volumes) ────────────────────────────

clean:
	$(COMPOSE) down -v
	@echo "All volumes removed."
