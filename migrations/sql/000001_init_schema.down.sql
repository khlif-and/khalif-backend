-- Down Migration: Init Schema

DROP INDEX IF EXISTS idx_admins_username;
DROP INDEX IF EXISTS idx_admins_email;
DROP TABLE IF EXISTS admin_audit_logs;
DROP TABLE IF EXISTS admins;
