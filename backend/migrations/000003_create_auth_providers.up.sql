CREATE TABLE auth_providers (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider     VARCHAR(20) NOT NULL,
    provider_uid VARCHAR(255) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(provider, provider_uid)
);

CREATE INDEX idx_auth_providers_user_id ON auth_providers(user_id);
CREATE INDEX idx_auth_providers_lookup ON auth_providers(provider, provider_uid);
