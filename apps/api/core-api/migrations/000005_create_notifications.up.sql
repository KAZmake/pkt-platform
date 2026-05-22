-- Уведомления пользователей: NATS consumer пишет, ЛК читает
CREATE TABLE IF NOT EXISTS notifications (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(50) NOT NULL CHECK (type IN ('payment', 'status', 'ticket', 'system')),
    title      VARCHAR(255) NOT NULL,
    body       TEXT        NOT NULL,
    is_read    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id      ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_is_read ON notifications (user_id, is_read);
