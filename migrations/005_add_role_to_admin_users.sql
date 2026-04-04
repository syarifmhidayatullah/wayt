-- +migrate Up
CREATE TYPE admin_role AS ENUM ('superadmin', 'admin');

ALTER TABLE admin_users
    ADD COLUMN role admin_role NOT NULL DEFAULT 'admin';

-- +migrate Down
ALTER TABLE admin_users DROP COLUMN role;
DROP TYPE IF EXISTS admin_role;
