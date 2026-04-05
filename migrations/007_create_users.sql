-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL       PRIMARY KEY,
    name       VARCHAR(100)    NOT NULL,
    phone      VARCHAR(20)     NOT NULL UNIQUE,
    password   VARCHAR(255)    NOT NULL,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Link booked queues to a user (nullable — QR-scanned queues have no user)
ALTER TABLE queues ADD COLUMN user_id BIGINT NULL REFERENCES users(id);

-- +migrate Down
ALTER TABLE queues DROP COLUMN IF EXISTS user_id;
DROP TABLE IF EXISTS users;
