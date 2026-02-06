CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    server_key TEXT,
    key_environment VARCHAR(50) NOT NULL DEFAULT 'production',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- indexes
CREATE INDEX idx_merchants_name ON merchants(name);