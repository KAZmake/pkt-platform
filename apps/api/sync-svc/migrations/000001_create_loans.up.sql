-- Займы из 1С: только sync-svc пишет, ручная правка запрещена
-- TimescaleDB hypertable по synced_at для временных запросов
CREATE TABLE IF NOT EXISTS loans (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    one_c_id    VARCHAR(255) NOT NULL,
    borrower_id UUID        NOT NULL,
    program_id  UUID,
    amount      NUMERIC(15, 2) NOT NULL,
    rate        NUMERIC(5, 2)  NOT NULL,
    term_months INT            NOT NULL,
    issued_at   DATE           NOT NULL,
    expires_at  DATE           NOT NULL,
    status      VARCHAR(50)    NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'overdue', 'closed')),
    synced_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    CONSTRAINT loans_one_c_id_key UNIQUE (one_c_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_loans_one_c_id    ON loans (one_c_id);
CREATE        INDEX IF NOT EXISTS idx_loans_borrower_id ON loans (borrower_id);

-- График платежей (из 1С)
CREATE TABLE IF NOT EXISTS payment_schedule (
    id        UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id   UUID        NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    due_date  DATE        NOT NULL,
    principal NUMERIC(15, 2) NOT NULL,
    interest  NUMERIC(15, 2) NOT NULL,
    total     NUMERIC(15, 2) NOT NULL,
    is_paid   BOOLEAN     NOT NULL DEFAULT FALSE,
    paid_at   DATE,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_schedule_loan_id  ON payment_schedule (loan_id);
-- BRIN эффективен для монотонно-возрастающих дат (sorted по времени)
CREATE INDEX IF NOT EXISTS idx_payment_schedule_due_date ON payment_schedule USING BRIN (due_date);

-- Задолженности (из 1С)
CREATE TABLE IF NOT EXISTS loan_debts (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id      UUID        NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    type         VARCHAR(50) NOT NULL CHECK (type IN ('principal', 'interest', 'penalty')),
    amount       NUMERIC(15, 2) NOT NULL,
    days_overdue INT            NOT NULL DEFAULT 0,
    synced_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_loan_debts_loan_id ON loan_debts (loan_id);
