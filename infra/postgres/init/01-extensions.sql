-- Задача 1.2: PostgreSQL 16 + PostGIS + TimescaleDB
-- Этот скрипт запускается автоматически при первом старте контейнера

-- Основная БД: расширения
\c pkt_db;

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;   -- полнотекстовый поиск
CREATE EXTENSION IF NOT EXISTS btree_gin; -- индексы GIN

SELECT postgis_version();
SELECT extversion FROM pg_extension WHERE extname = 'timescaledb';
