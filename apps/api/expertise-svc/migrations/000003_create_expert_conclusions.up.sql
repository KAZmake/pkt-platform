-- Заключения экспертов по трём этапам экспертизы + голосование КК
CREATE TABLE IF NOT EXISTS expert_conclusions (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id  UUID        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    expert_id       UUID        NOT NULL,
    stage           VARCHAR(50) NOT NULL
                    CHECK (stage IN ('collateral_expertise', 'legal_check', 'credit_analysis')),
    risks           JSONB,
    conclusion_text TEXT,
    result          VARCHAR(50) NOT NULL CHECK (result IN ('approved', 'rejected', 'revision')),
    file_path       VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_expert_conclusions_application_id ON expert_conclusions (application_id);

-- Голосование членов кредитного комитета
CREATE TABLE IF NOT EXISTS committee_votes (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    expert_id      UUID        NOT NULL,
    vote           VARCHAR(50) NOT NULL CHECK (vote IN ('approved', 'rejected', 'abstained')),
    comment        TEXT,
    signed_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_committee_votes_application_id ON committee_votes (application_id);
