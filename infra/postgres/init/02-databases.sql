-- Создаём отдельные БД для Keycloak и Directus
-- Запускается от имени суперпользователя postgres

-- Keycloak DB
CREATE USER keycloak WITH PASSWORD 'keycloak_secret';
CREATE DATABASE keycloak_db OWNER keycloak ENCODING 'UTF8';
GRANT ALL PRIVILEGES ON DATABASE keycloak_db TO keycloak;

-- Directus DB
CREATE USER directus WITH PASSWORD 'directus_secret';
CREATE DATABASE directus_db OWNER directus ENCODING 'UTF8';
GRANT ALL PRIVILEGES ON DATABASE directus_db TO directus;

-- Подключаемся к directus_db и включаем uuid-ossp (нужен Directus)
\c directus_db;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
