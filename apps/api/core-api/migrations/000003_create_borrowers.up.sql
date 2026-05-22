-- Заёмщики: расширенный профиль пользователя с ИИН/БИН
CREATE TABLE IF NOT EXISTS borrowers (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    inn           VARCHAR(50) NOT NULL,
    bin           VARCHAR(50),
    org_name      VARCHAR(255),
    activity_type VARCHAR(50) CHECK (activity_type IN ('crop_farming', 'livestock', 'mixed')),
    farm_id       UUID        REFERENCES farms(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT borrowers_inn_key UNIQUE (inn)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_borrowers_inn     ON borrowers (inn);
CREATE        INDEX IF NOT EXISTS idx_borrowers_user_id ON borrowers (user_id);
