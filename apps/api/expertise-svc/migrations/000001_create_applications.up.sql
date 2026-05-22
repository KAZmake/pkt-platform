-- Заявки на кредит + FSM-история (только INSERT, никогда UPDATE/DELETE)
CREATE TABLE IF NOT EXISTS applications (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id  UUID        NOT NULL,
    program_id   UUID        NOT NULL,
    assignee_id  UUID,
    status       VARCHAR(50) NOT NULL DEFAULT 'received'
                 CHECK (status IN (
                     'received', 'primary_scoring', 'security_check',
                     'collateral_expertise', 'legal_check', 'credit_analysis',
                     'credit_committee', 'approved', 'rejected', 'revision',
                     'documentation', 'issued'
                 )),
    amount       NUMERIC(15, 2) NOT NULL,
    term_months  INT            NOT NULL,
    payment_type VARCHAR(50)    NOT NULL CHECK (payment_type IN ('annuity', 'differentiated')),
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_borrower_id ON applications (borrower_id);
CREATE INDEX IF NOT EXISTS idx_applications_status      ON applications (status);

-- Аудит-лог переходов FSM. Правило: только INSERT, UPDATE и DELETE запрещены на уровне БД.
CREATE TABLE IF NOT EXISTS application_history (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID        NOT NULL REFERENCES applications(id),
    from_status    VARCHAR(50),
    to_status      VARCHAR(50) NOT NULL,
    actor_id       UUID        NOT NULL,
    comment        TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Защита аудит-лога от модификации
CREATE OR REPLACE RULE application_history_no_update
    AS ON UPDATE TO application_history DO INSTEAD NOTHING;

CREATE OR REPLACE RULE application_history_no_delete
    AS ON DELETE TO application_history DO INSTEAD NOTHING;

CREATE INDEX IF NOT EXISTS idx_application_history_application_id ON application_history (application_id);
