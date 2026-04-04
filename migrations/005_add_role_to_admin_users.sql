-- +migrate Up
ALTER TABLE admin_users
    ADD COLUMN role ENUM('superadmin', 'admin') NOT NULL DEFAULT 'admin' AFTER username;

-- +migrate Down
ALTER TABLE admin_users DROP COLUMN role;
