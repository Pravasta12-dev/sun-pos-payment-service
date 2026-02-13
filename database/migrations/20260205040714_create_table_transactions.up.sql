CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id UUID,
    order_id VARCHAR(100) NOT NULL UNIQUE,
    amount DECIMAL(15, 2) NOT NULL,
    payment_type VARCHAR(50) NOT NULL DEFAULT 'qris',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    bill_id VARCHAR(100) NULL,

    paid_at TIMESTAMP NULL,
    qr_url TEXT NULL,
    expired_at TIMESTAMP NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_transactions_merchant
        FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

CREATE INDEX idx_transactions_merchant_id ON transactions(merchant_id);
CREATE INDEX idx_transactions_order_id ON transactions(order_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_expired_at ON transactions(expired_at);
CREATE INDEX idx_transactions_payment_type ON transactions(payment_type);
CREATE INDEX idx_transactions_bill_id ON transactions(bill_id);