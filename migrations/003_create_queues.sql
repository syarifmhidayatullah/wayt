-- +migrate Up
CREATE TABLE IF NOT EXISTS queues (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    branch_id    BIGINT UNSIGNED NOT NULL,
    qr_token     VARCHAR(36)     NOT NULL,
    queue_number VARCHAR(20)     NOT NULL,
    status       ENUM('waiting','called','done','expired') NOT NULL DEFAULT 'waiting',
    created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_queues_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_queues_branch_id ON queues(branch_id);
CREATE INDEX idx_queues_qr_token  ON queues(qr_token);
CREATE INDEX idx_queues_status    ON queues(status);

-- +migrate Down
DROP TABLE IF EXISTS queues;
