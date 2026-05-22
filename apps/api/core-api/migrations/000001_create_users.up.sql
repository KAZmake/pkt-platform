-- Users: зеркало Keycloak, источник истины по ролям
CREATE TABLE IF NOT EXISTS users (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    keycloak_id VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    role        VARCHAR(50)  NOT NULL DEFAULT 'borrower'
                CHECK (role IN ('public', 'borrower', 'employee', 'expert', 'admin')),
    first_name  VARCHAR(255),
    last_name   VARCHAR(255),
    phone       VARCHAR(50),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_keycloak_id_key UNIQUE (keycloak_id),
    CONSTRAINT users_email_key       UNIQUE (email)
);

-- Критичный индекс: каждый API-запрос проверяет JWT sub → keycloak_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_keycloak_id ON users (keycloak_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email       ON users (email);
