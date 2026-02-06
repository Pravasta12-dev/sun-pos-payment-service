CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id UUID NOT NULL,
    order_id VARCHAR(100) NOT NULL UNIQUE,
    amount DECIMAL(15, 2) NOT NULL,
    payment_type VARCHAR(50) NOT NULL DEFAULT 'qris',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',

    paid_at TIMESTAMP NULL,
    expired_at TIMESTAMP NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_transactions_merchant
        FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);