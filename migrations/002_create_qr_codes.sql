-- +migrate Up
CREATE TABLE IF NOT EXISTS qr_codes (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    branch_id   BIGINT UNSIGNED NOT NULL,
    token       VARCHAR(36)     NOT NULL UNIQUE,
    is_active   TINYINT(1)      NOT NULL DEFAULT 1,
    expired_at  DATETIME        NOT NULL,
    created_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_qr_codes_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_qr_codes_token ON qr_codes(token);
CREATE INDEX idx_qr_codes_branch_id ON qr_codes(branch_id);

-- +migrate Down
DROP TABLE IF EXISTS qr_codes;
