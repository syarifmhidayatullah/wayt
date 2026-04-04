-- +migrate Up
CREATE TABLE IF NOT EXISTS qr_codes (
    id         BIGSERIAL       PRIMARY KEY,
    branch_id  BIGINT          NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    token      VARCHAR(36)     NOT NULL UNIQUE,
    is_active  BOOLEAN         NOT NULL DEFAULT TRUE,
    expired_at TIMESTAMPTZ     NOT NULL,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qr_codes_token     ON qr_codes(token);
CREATE INDEX idx_qr_codes_branch_id ON qr_codes(branch_id);

-- +migrate Down
DROP TABLE IF EXISTS qr_codes;
