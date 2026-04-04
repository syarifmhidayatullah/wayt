-- +migrate Up
CREATE TABLE IF NOT EXISTS admin_users (
    id         BIGSERIAL       PRIMARY KEY,
    username   VARCHAR(100)    NOT NULL UNIQUE,
    password   VARCHAR(255)    NOT NULL,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS admin_users;
