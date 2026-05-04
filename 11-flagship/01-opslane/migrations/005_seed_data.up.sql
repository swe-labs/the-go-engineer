-- 005_seed_data.up.sql
-- Seed data for development

-- Insert a demo tenant
INSERT INTO tenants (id, name, slug, created_at) 
VALUES (1, 'Demo Organization', 'demo', NOW())
ON CONFLICT (slug) DO NOTHING;

-- Insert demo users (password is "password123" hashed with bcrypt)
-- $2a$10$xCJ7LzR6N5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X5X
INSERT INTO users (id, tenant_id, email, display_name, password_hash, role, created_at) 
VALUES 
    (1, 1, 'admin@demo.com', 'Demo Admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJ8P2P8O8O8O8O8O8O8O8O8O8O', 'admin', NOW()),
    (2, 1, 'user@demo.com', 'Demo User', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJ8P2P8O8O8O8O8O8O8O8O8O8O', 'user', NOW())
ON CONFLICT (tenant_id, email) DO NOTHING;

-- Insert demo orders
INSERT INTO orders (id, tenant_id, user_id, status, total_cents, currency, idempotency_key, created_at, updated_at)
VALUES 
    (1, 1, 1, 'completed', 10000, 'USD', 'order-001', NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
    (2, 1, 1, 'pending', 25000, 'USD', 'order-002', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
    (3, 1, 2, 'completed', 5000, 'USD', 'order-003', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day')
ON CONFLICT (tenant_id, idempotency_key) DO NOTHING;

-- Insert demo payments
INSERT INTO payments (id, tenant_id, order_id, status, provider_reference, amount_cents, failure_reason, created_at, updated_at)
VALUES 
    (1, 1, 1, 'succeeded', 'pi_001', 10000, '', NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
    (2, 1, 2, 'pending', 'pi_002', 25000, '', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
    (3, 1, 3, 'succeeded', 'pi_003', 5000, '', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day')
ON CONFLICT (tenant_id, provider_reference) DO NOTHING;