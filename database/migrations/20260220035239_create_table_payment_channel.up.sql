CREATE TABLE IF NOT EXISTS payment_channels (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL,
    label VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    fee_type VARCHAR(20) NOT NULL DEFAULT 'percentage',
    fee_value DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_channels_type ON payment_channels(type);
