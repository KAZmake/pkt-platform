-- Программы кредитования: публичный справочник, управляется через Directus CMS
CREATE TABLE IF NOT EXISTS loan_programs (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    name_kz         VARCHAR(255),
    name_en         VARCHAR(255),
    rate            NUMERIC(5, 2) NOT NULL,
    min_amount      NUMERIC(15, 2) NOT NULL,
    max_amount      NUMERIC(15, 2) NOT NULL,
    min_term_months INT          NOT NULL,
    max_term_months INT          NOT NULL,
    activity_types  TEXT[]       NOT NULL DEFAULT '{}',
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_loan_programs_is_active ON loan_programs (is_active);
