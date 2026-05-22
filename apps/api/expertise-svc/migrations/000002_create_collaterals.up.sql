-- Залоги: самостоятельные сущности, могут переходить между кредитами
CREATE TABLE IF NOT EXISTS collaterals (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    type                VARCHAR(50) NOT NULL
                        CHECK (type IN ('land', 'equipment', 'livestock', 'real_estate', 'other')),
    description         TEXT,
    estimated_value     NUMERIC(15, 2),
    cadastral_number    VARCHAR(255),
    insurance_expiry    DATE,
    last_inventory_date DATE,
    is_released         BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Связь M:N заявка ↔ залог с датами привязки/высвобождения
CREATE TABLE IF NOT EXISTS application_collaterals (
    application_id UUID        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    collateral_id  UUID        NOT NULL REFERENCES collaterals(id)  ON DELETE CASCADE,
    attached_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at    TIMESTAMPTZ,
    PRIMARY KEY (application_id, collateral_id)
);
