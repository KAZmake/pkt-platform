-- Обращения заёмщиков (тикеты) + сообщения
CREATE TABLE IF NOT EXISTS tickets (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    borrower_id UUID        NOT NULL REFERENCES borrowers(id) ON DELETE CASCADE,
    assignee_id UUID        REFERENCES users(id),
    type        VARCHAR(50) NOT NULL CHECK (type IN ('early_repayment', 'restructuring', 'prolongation', 'other')),
    subject     VARCHAR(500) NOT NULL,
    status      VARCHAR(50) NOT NULL DEFAULT 'open'
                CHECK (status IN ('open', 'in_progress', 'resolved', 'closed')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tickets_borrower_id ON tickets (borrower_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status      ON tickets (status);

CREATE TABLE IF NOT EXISTS ticket_messages (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id       UUID        NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    author_id       UUID        NOT NULL REFERENCES users(id),
    body            TEXT        NOT NULL,
    attachment_path VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ticket_messages_ticket_id ON ticket_messages (ticket_id);
