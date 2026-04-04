-- +migrate Up
CREATE TABLE IF NOT EXISTS branches (
    id             BIGSERIAL       PRIMARY KEY,
    name           VARCHAR(100)    NOT NULL,
    prefix         VARCHAR(10)     NOT NULL,
    is_active      BOOLEAN         NOT NULL DEFAULT TRUE,
    current_number INTEGER         NOT NULL DEFAULT 0,
    last_number    INTEGER         NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ     NULL
);

-- +migrate Down
DROP TABLE IF EXISTS branches;
