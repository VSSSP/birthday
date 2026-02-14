CREATE TABLE recipients (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    age        INT NOT NULL DEFAULT 0,
    gender     VARCHAR(20) NOT NULL DEFAULT 'other',
    min_budget DECIMAL(10, 2) NOT NULL DEFAULT 0,
    max_budget DECIMAL(10, 2) NOT NULL DEFAULT 0,
    keywords   TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recipients_user_id ON recipients(user_id);
