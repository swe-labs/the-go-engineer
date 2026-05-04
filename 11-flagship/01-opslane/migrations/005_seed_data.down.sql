-- 005_seed_data.down.sql
-- Remove seed data

DELETE FROM payments WHERE tenant_id = 1;
DELETE FROM orders WHERE tenant_id = 1;
DELETE FROM users WHERE tenant_id = 1;
DELETE FROM tenants WHERE id = 1;