-- Документы: полиморфная привязка к заёмщику / заявке / залогу
-- Файлы хранятся физически в MinIO, здесь только метаданные
CREATE TABLE IF NOT EXISTS documents (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id   UUID        NOT NULL,
    owner_type VARCHAR(50) NOT NULL CHECK (owner_type IN ('borrower', 'application', 'collateral')),
    bucket     VARCHAR(255) NOT NULL,
    object_key VARCHAR(500) NOT NULL,
    name       VARCHAR(500) NOT NULL,
    mime_type  VARCHAR(100),
    size_bytes BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_documents_owner ON documents (owner_id, owner_type);
