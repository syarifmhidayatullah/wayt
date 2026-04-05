-- +migrate Up

-- Create counters table (loket/service point per branch, with own prefix & numbering)
CREATE TABLE IF NOT EXISTS counters (
    id             BIGSERIAL       PRIMARY KEY,
    branch_id      BIGINT          NOT NULL REFERENCES branches(id),
    name           VARCHAR(100)    NOT NULL,
    prefix         VARCHAR(10)     NOT NULL,
    is_active      BOOLEAN         NOT NULL DEFAULT TRUE,
    current_number INTEGER         NOT NULL DEFAULT 0,
    last_number    INTEGER         NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ     NULL
);

-- Migrate existing branch prefix/numbering data into counters (1 counter per branch)
INSERT INTO counters (branch_id, name, prefix, is_active, current_number, last_number, created_at, updated_at)
SELECT id, name, prefix, is_active, current_number, last_number, created_at, updated_at
FROM branches
WHERE deleted_at IS NULL;

-- Add counter_id to qr_codes
ALTER TABLE qr_codes ADD COLUMN counter_id BIGINT NULL REFERENCES counters(id);

-- Fill counter_id in qr_codes from branch → migrated counter
UPDATE qr_codes qr
SET counter_id = c.id
FROM counters c
WHERE c.branch_id = qr.branch_id;

-- Make counter_id NOT NULL after fill
ALTER TABLE qr_codes ALTER COLUMN counter_id SET NOT NULL;

-- Add counter_id to queues
ALTER TABLE queues ADD COLUMN counter_id BIGINT NULL REFERENCES counters(id);

-- Fill counter_id in queues from branch → migrated counter
UPDATE queues q
SET counter_id = c.id
FROM counters c
WHERE c.branch_id = q.branch_id;

-- Make counter_id NOT NULL after fill
ALTER TABLE queues ALTER COLUMN counter_id SET NOT NULL;

-- Add branch_id to admin_users (nullable — only set for role='admin')
ALTER TABLE admin_users ADD COLUMN branch_id BIGINT NULL REFERENCES branches(id);

-- Remove obsolete columns from branches
ALTER TABLE branches DROP COLUMN IF EXISTS prefix;
ALTER TABLE branches DROP COLUMN IF EXISTS current_number;
ALTER TABLE branches DROP COLUMN IF EXISTS last_number;

-- +migrate Down
ALTER TABLE branches ADD COLUMN prefix VARCHAR(10) NOT NULL DEFAULT '';
ALTER TABLE branches ADD COLUMN current_number INTEGER NOT NULL DEFAULT 0;
ALTER TABLE branches ADD COLUMN last_number INTEGER NOT NULL DEFAULT 0;
ALTER TABLE admin_users DROP COLUMN IF EXISTS branch_id;
ALTER TABLE queues DROP COLUMN IF EXISTS counter_id;
ALTER TABLE qr_codes DROP COLUMN IF EXISTS counter_id;
DROP TABLE IF EXISTS counters;
