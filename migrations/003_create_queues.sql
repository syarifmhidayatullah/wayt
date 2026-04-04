-- +migrate Up
CREATE TYPE queue_status AS ENUM ('waiting', 'called', 'done', 'expired');

CREATE TABLE IF NOT EXISTS queues (
    id           BIGSERIAL       PRIMARY KEY,
    branch_id    BIGINT          NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    qr_token     VARCHAR(36)     NOT NULL,
    queue_number VARCHAR(20)     NOT NULL,
    status       queue_status    NOT NULL DEFAULT 'waiting',
    created_at   TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_queues_branch_id ON queues(branch_id);
CREATE INDEX idx_queues_qr_token  ON queues(qr_token);
CREATE INDEX idx_queues_status    ON queues(status);

-- +migrate Down
DROP TABLE IF EXISTS queues;
DROP TYPE IF EXISTS queue_status;
